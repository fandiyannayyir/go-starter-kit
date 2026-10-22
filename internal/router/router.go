// Package router is the single place that mounts every handler's route table
// onto Echo. Adding a domain means accepting its handler here and adding one
// reg.Add line — the routing surface stays in one file, out of cmd.
package router

import (
	"github.com/labstack/echo/v4"

	"go-starter-kit/internal/handler"
	"go-starter-kit/pkg/route"
)

// Build registers all application routes on e, returning an error on a
// duplicate method+path across any handler.
func Build(e *echo.Echo, health *handler.HealthHandler, example *handler.ExampleHandler) error {
	reg := route.NewRegistry()
	reg.Add(health.Routes()...)
	reg.Add(example.Routes()...)
	return reg.Register(e)
}
