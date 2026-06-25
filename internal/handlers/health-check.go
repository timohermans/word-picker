package handlers

import (
	"hert/gotest/internal/db"
	"net/http"

	"github.com/labstack/echo/v5"
)

const urlHealthCheck = "/health"
const urlReadyCheck = "/ready"

func RegisterHealthCheck(e *echo.Echo, queries *db.Queries) {
	e.GET(urlHealthCheck, func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET(urlReadyCheck, func(c *echo.Context) error {
		ctx := c.Request().Context()

		if isReady, err := queries.IsReady(ctx); err != nil || !isReady {
			return c.NoContent(http.StatusServiceUnavailable)
		}

		return c.NoContent(http.StatusOK)
	})
}
