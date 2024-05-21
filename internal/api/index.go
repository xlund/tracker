package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xlund/tracker/internal/view/layout"
	"github.com/xlund/tracker/internal/view/page"
)

func (a *api) getIndex(c echo.Context) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	users, err := a.userRepo.GetAll(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	r := layout.Base("Chess Tournament Tracker", "/", page.Index(users))
	return r.Render(ctx, c.Response().Writer)
}
