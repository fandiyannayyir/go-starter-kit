package handler

import (
	"net/http"

	"go-starter-kit/pkg/route"
)

// Routes declares the health URL surface (public, no prefix).
func (h *HealthHandler) Routes() []route.RouteGroup {
	return []route.RouteGroup{{
		Routes: []route.Route{
			{Method: http.MethodGet, Path: "/health", Handler: h.Live, Public: true},
			{Method: http.MethodGet, Path: "/ready", Handler: h.Ready, Public: true},
		},
	}}
}
