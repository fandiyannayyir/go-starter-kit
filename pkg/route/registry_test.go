package route

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func noop(c echo.Context) error { return nil }

func TestRegister_DuplicateIsRejected(t *testing.T) {
	reg := NewRegistry()
	reg.Add(RouteGroup{
		Prefix: "/v1",
		Routes: []Route{
			{Method: "GET", Path: "/ping", Handler: noop},
			{Method: "GET", Path: "/ping", Handler: noop}, // same method+path
		},
	})
	if err := reg.Register(echo.New()); err == nil {
		t.Fatal("expected duplicate route error, got nil")
	}
}

func TestRegister_DistinctRoutesOK(t *testing.T) {
	reg := NewRegistry()
	reg.Add(RouteGroup{
		Prefix: "/v1",
		Routes: []Route{
			{Method: "GET", Path: "/ping", Handler: noop},
			{Method: "POST", Path: "/ping", Handler: noop},
			{Method: "GET", Path: "/health", Handler: noop},
		},
	})
	if err := reg.Register(echo.New()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
