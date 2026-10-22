package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"go-starter-kit/internal/shared/log"
)

// Logger logs one structured line per request: method, path, status, latency
// and request id. It runs closest to the handler so the latency it records is
// the handler's, and it logs via defer so a panic unwinding through it is still
// recorded before Recover converts it to a 500.
func Logger(logger log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			start := time.Now()
			req := c.Request()
			defer func() {
				fields := map[string]any{
					"method":     req.Method,
					"path":       req.URL.Path,
					"status":     c.Response().Status,
					"latency_ms": time.Since(start).Milliseconds(),
					"request_id": req.Header.Get(echo.HeaderXRequestID),
				}
				entry := logger.WithContext(req.Context()).WithFields(fields)
				if err != nil {
					entry.WithError(err).Error("request failed")
					return
				}
				entry.Info("request")
			}()
			err = next(c)
			return err
		}
	}
}
