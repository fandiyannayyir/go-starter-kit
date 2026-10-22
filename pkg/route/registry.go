// Package route provides a declarative route table for Echo.
//
// Routes are declared as data in one place per handler module and registered
// centrally, so every endpoint is traceable and duplicate method+path pairs
// are caught at startup instead of silently shadowing each other.
package route

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Route is a single route entry.
type Route struct {
	Method     string
	Path       string
	Handler    echo.HandlerFunc
	Middleware []echo.MiddlewareFunc
	Public     bool // if true, auth middleware skips this route
}

// RouteGroup is a set of routes sharing a prefix and middleware.
type RouteGroup struct {
	Prefix     string
	Middleware []echo.MiddlewareFunc
	Routes     []Route
}

// Registry collects route groups before they are applied to Echo.
type Registry struct {
	groups []RouteGroup
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry { return &Registry{} }

// Add appends one or more route groups.
func (r *Registry) Add(groups ...RouteGroup) { r.groups = append(r.groups, groups...) }

// Register applies all routes to e, returning an error on a duplicate
// method+path so the collision is caught at boot instead of at runtime.
func (r *Registry) Register(e *echo.Echo) error {
	seen := make(map[string]bool)
	for _, g := range r.groups {
		group := e.Group(g.Prefix, g.Middleware...)
		for _, rt := range g.Routes {
			key := rt.Method + " " + g.Prefix + rt.Path
			if seen[key] {
				return fmt.Errorf("route: duplicate %s", key)
			}
			seen[key] = true
			group.Add(rt.Method, rt.Path, rt.Handler, rt.Middleware...)
		}
	}
	return nil
}
