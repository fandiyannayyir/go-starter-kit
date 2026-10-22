package handler

import (
	"net/http"

	"go-starter-kit/pkg/route"
)

// Routes declares the example URL surface as data, kept separate from the
// handler methods so the route table is readable in one place.
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
