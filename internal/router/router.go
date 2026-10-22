// Package router is the single place that mounts every handler's routes onto
// Echo. One function per domain keeps it modular; Build aggregates them. Adding
// a domain = add its <domain>Routes func and one reg.Add line in Build.
package router

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-starter-kit/internal/handler"
	"go-starter-kit/pkg/route"
)

// Build registers all application routes on e, returning an error on a
// duplicate method+path across any domain.
func Build(e *echo.Echo, health *handler.HealthHandler, example *handler.ExampleHandler) error {
	reg := route.NewRegistry()
	reg.Add(healthRoutes(health)...)
	reg.Add(exampleRoutes(example)...)
	return reg.Register(e)
}

func healthRoutes(h *handler.HealthHandler) []route.RouteGroup {
	return []route.RouteGroup{{
		Routes: []route.Route{
			{Method: http.MethodGet, Path: "/health", Handler: h.Live, Public: true},
			{Method: http.MethodGet, Path: "/ready", Handler: h.Ready, Public: true},
		},
	}}
}

func exampleRoutes(h *handler.ExampleHandler) []route.RouteGroup {
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
