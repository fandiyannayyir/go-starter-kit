// Package server builds the configured Echo instance (middleware chain +
// centralized error handler) and runs it with graceful shutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"go-starter-kit/internal/middleware"
	"go-starter-kit/internal/shared/log"
	"go-starter-kit/pkg/apperror"
	"go-starter-kit/pkg/response"
)

// shutdownTimeout bounds how long in-flight requests have to drain on shutdown.
const shutdownTimeout = 10 * time.Second

// New returns an Echo instance with the centralized error handler and the base
// middleware chain applied, in the order:
//
//	RequestID → CORS → Recover → Logger → handler
//
// (JWT/RBAC, Gzip and OTel are added per the slice that needs them.)
func New(logger log.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = errorHandler(logger)

	e.Use(middleware.RequestID())
	e.Use(echomw.CORS())
	e.Use(echomw.Recover())
	e.Use(middleware.Logger(logger))

	return e
}

// errorHandler maps any error returned by a handler to a response.Envelope:
//   - apperror.AppError  → its own code + HTTP status
//   - *echo.HTTPError    → HTTP_<status> + its message (routing/binding errors)
//   - anything else      → 500 with a generic message (details only in the log)
func errorHandler(logger log.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		r := c.Request()

		var appErr apperror.AppError
		if errors.As(err, &appErr) {
			_ = c.JSON(appErr.HTTPStatus(), response.Error(r, appErr.Code(), appErr.Error()))
			return
		}

		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			_ = c.JSON(httpErr.Code, response.Error(r,
				fmt.Sprintf("HTTP_%d", httpErr.Code),
				fmt.Sprint(httpErr.Message)))
			return
		}

		logger.WithContext(r.Context()).WithError(err).Error("unhandled error")
		_ = c.JSON(http.StatusInternalServerError,
			response.Error(r, apperror.CodeInternal, "internal server error"))
	}
}

// Start runs e on addr and blocks until an interrupt/SIGTERM, then shuts down
// gracefully within shutdownTimeout.
func Start(e *echo.Echo, addr string, logger log.Logger) error {
	go func() {
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal("http server stopped unexpectedly")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return e.Shutdown(ctx)
}
