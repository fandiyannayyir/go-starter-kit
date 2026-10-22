// Package middleware holds the project's custom Echo middleware. CORS and
// Recover use Echo's built-ins directly; only the pieces with project-specific
// behavior (request id propagation, structured request logging) live here.
package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestID ensures every request carries an X-Request-ID header. If the client
// did not send one, it generates a UUID and sets it on both the request and the
// response, so the id logged, returned in the response envelope, and seen by the
// client all match.
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = uuid.NewString()
				req.Header.Set(echo.HeaderXRequestID, id)
			}
			c.Response().Header().Set(echo.HeaderXRequestID, id)
			return next(c)
		}
	}
}
