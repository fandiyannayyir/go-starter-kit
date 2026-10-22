// Package handler holds the HTTP handlers and their declarative route tables.
package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"go-starter-kit/configs/db"
	"go-starter-kit/internal/shared/cache"
	"go-starter-kit/pkg/response"
)

// HealthHandler reports process liveness and dependency readiness.
type HealthHandler struct {
	db    *db.Database
	cache cache.Cache
}

func NewHealthHandler(database *db.Database, c cache.Cache) *HealthHandler {
	return &HealthHandler{db: database, cache: c}
}

// Live answers whether the process is up. It touches no dependency, so a load
// balancer can tell "process running" apart from "dependencies healthy".
func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK,
		response.OK(c.Request(), map[string]string{"status": "ok"}, "alive"))
}

// Ready pings the database and cache. It returns 503 if either is unreachable,
// so orchestrators hold traffic until the service can actually serve it.
func (h *HealthHandler) Ready(c echo.Context) error {
	r := c.Request()
	checks := map[string]string{"db": "up", "cache": "up"}
	ready := true

	if err := h.db.PingDb(); err != nil {
		checks["db"] = "down"
		ready = false
	}
	if err := h.cache.Ping(r.Context()); err != nil {
		checks["cache"] = "down"
		ready = false
	}

	if !ready {
		return c.JSON(http.StatusServiceUnavailable, response.Error(r, "NOT_READY",
			fmt.Sprintf("dependency unavailable (db=%s cache=%s)", checks["db"], checks["cache"])))
	}
	return c.JSON(http.StatusOK, response.OK(r, checks, "ready"))
}
