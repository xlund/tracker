package api

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/view/layout"
	"github.com/xlund/tracker/internal/view/page"
)

func (a *api) createUser(e echo.Context) error {

	_, cancel := context.WithCancel(e.Request().Context())
	defer cancel()
	return e.String(http.StatusOK, "create user")

}

func (a *api) getUser(c echo.Context) error {
	session, err := session.Get("auth-session", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if session.Values["profile"] == nil {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	profile := domain.UserAUthProfileFromMap(session.Values["profile"].(map[string]interface{}))

	t := layout.Base("User", page.UserProfile(profile))

	return t.Render(c.Request().Context(), c.Response().Writer)
}

func (a *api) getUsers(e echo.Context) error {
	ctx, cancel := context.WithCancel(e.Request().Context())
	defer cancel()

	users, err := a.userRepo.GetAll(ctx)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	c := layout.Base("Users", page.Users(users))
	return c.Render(ctx, e.Response().Writer)

}

func (a *api) deleteUser(e echo.Context) error {
	ctx, cancel := context.WithCancel(e.Request().Context())
	defer cancel()

	id := e.Param("id")
	log.Default().Printf("deleting user %s", id)
	err := a.userRepo.Delete(ctx, id)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "deleted user")
}
