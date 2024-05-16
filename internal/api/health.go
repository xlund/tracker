package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *api) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
