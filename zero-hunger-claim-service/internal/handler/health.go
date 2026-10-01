package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Health is a liveness probe: it only reports that the process is running.
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// Ready is a readiness probe: it runs check (for example a database ping) and
// returns 503 when the dependency is unavailable.
func Ready(check func(ctx context.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()

		if err := check(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"error":  err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
	}
}
