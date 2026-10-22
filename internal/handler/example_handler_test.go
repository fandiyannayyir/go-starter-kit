package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	domain "go-starter-kit/internal/domain/example"
	"go-starter-kit/internal/handler"
	"go-starter-kit/internal/router"
	"go-starter-kit/internal/server"
	example "go-starter-kit/internal/service/example"
	"go-starter-kit/internal/shared/log"
)

// fakeService implements example.Service in-memory, so these tests exercise the
// real HTTP chain (middleware -> validator -> router -> handler -> central error
// handler -> envelope) without a database.
type fakeService struct{}

func (f *fakeService) Create(_ context.Context, req example.CreateRequest) (*example.Response, error) {
	return &example.Response{ID: "fake-id", Name: req.Name, Description: req.Description}, nil
}

func (f *fakeService) Get(_ context.Context, _ string) (*example.Response, error) {
	return nil, domain.ErrNotFound
}

func (f *fakeService) List(_ context.Context, _, _ int) ([]example.Response, int, error) {
	return []example.Response{}, 0, nil
}

func (f *fakeService) Update(_ context.Context, _ string, _ example.UpdateRequest) (*example.Response, error) {
	return nil, domain.ErrNotFound
}

func (f *fakeService) Delete(_ context.Context, _ string) error { return nil }

func newTestServer(t *testing.T) *echo.Echo {
	t.Helper()
	e := server.New(log.GetLogger())
	// Health handler gets nil deps: only its route table is mounted here and
	// these tests never call /ready (the only endpoint that touches them).
	health := handler.NewHealthHandler(nil, nil)
	if err := router.Build(e, health, handler.NewExampleHandler(&fakeService{})); err != nil {
		t.Fatalf("build routes: %v", err)
	}
	return e
}

func do(t *testing.T, e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

type envelope struct {
	Success bool `json:"success"`
	Error   *struct {
		Code   string `json:"code"`
		Fields []struct {
			Field string `json:"field"`
		} `json:"fields"`
	} `json:"error"`
	RequestID string `json:"request_id"`
}

func parse(t *testing.T, rec *httptest.ResponseRecorder) envelope {
	t.Helper()
	var env envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("parse envelope: %v (body: %s)", err, rec.Body.String())
	}
	return env
}

func TestCreate_Created(t *testing.T) {
	rec := do(t, newTestServer(t), http.MethodPost, "/v1/examples", `{"name":"produk"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if env := parse(t, rec); !env.Success || env.RequestID == "" {
		t.Fatalf("want success envelope with request_id, got %s", rec.Body.String())
	}
}

func TestCreate_ValidationErrorEnvelope(t *testing.T) {
	rec := do(t, newTestServer(t), http.MethodPost, "/v1/examples", `{"description":"tanpa nama"}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	env := parse(t, rec)
	if env.Error == nil || env.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("want VALIDATION_ERROR, got %s", rec.Body.String())
	}
	if len(env.Error.Fields) != 1 || env.Error.Fields[0].Field != "name" {
		t.Fatalf("want field error on 'name', got %s", rec.Body.String())
	}
}

func TestGet_DomainErrorMapsToEnvelope(t *testing.T) {
	rec := do(t, newTestServer(t), http.MethodGet, "/v1/examples/missing", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if env := parse(t, rec); env.Error == nil || env.Error.Code != "EXAMPLE_NOT_FOUND" {
		t.Fatalf("want EXAMPLE_NOT_FOUND, got %s", rec.Body.String())
	}
}
