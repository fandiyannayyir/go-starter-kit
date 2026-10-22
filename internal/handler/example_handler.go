package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go-starter-kit/internal/service"
	"go-starter-kit/pkg/response"
	"go-starter-kit/pkg/route"
)

// ExampleHandler is the HTTP surface of the example domain. Handlers return the
// error from the service as-is; the centralized error handler maps AppError to
// the response envelope and the correct status.
type ExampleHandler struct {
	svc service.ExampleService
}

func NewExampleHandler(svc service.ExampleService) *ExampleHandler {
	return &ExampleHandler{svc: svc}
}

func (h *ExampleHandler) Create(c echo.Context) error {
	var req service.CreateExampleRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	resp, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response.Created(c.Request(), resp, "example created"))
}

func (h *ExampleHandler) Get(c echo.Context) error {
	resp, err := h.svc.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.OK(c.Request(), resp, "ok"))
}

func (h *ExampleHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	items, total, err := h.svc.List(c.Request().Context(), (page-1)*perPage, perPage)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.Paginated(c.Request(), items, page, perPage, total, "ok"))
}

func (h *ExampleHandler) Update(c echo.Context) error {
	var req service.UpdateExampleRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	resp, err := h.svc.Update(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.OK(c.Request(), resp, "example updated"))
}

func (h *ExampleHandler) Delete(c echo.Context) error {
	if err := h.svc.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response.OK[any](c.Request(), nil, "example deleted"))
}

// Routes returns the example route groups (mounted by cmd into the registry).
func (h *ExampleHandler) Routes() []route.RouteGroup {
	return []route.RouteGroup{{
		Prefix: "/v1/examples",
		Routes: []route.Route{
			{Method: http.MethodPost, Path: "", Handler: h.Create},
			{Method: http.MethodGet, Path: "", Handler: h.List},
			{Method: http.MethodGet, Path: "/:id", Handler: h.Get},
			{Method: http.MethodPut, Path: "/:id", Handler: h.Update},
			{Method: http.MethodDelete, Path: "/:id", Handler: h.Delete},
		},
	}}
}
