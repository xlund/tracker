package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/view/layout"
	"github.com/xlund/tracker/internal/view/page"
	"github.com/xlund/tracker/internal/view/partial/table/item"
)

func (a *api) createGame(e echo.Context) error {
	ctx, cancel := context.WithCancel(e.Request().Context())
	defer cancel()

	whiteID := e.FormValue("white")
	blackID := e.FormValue("black")
	if whiteID == blackID {
		return e.String(http.StatusBadRequest, "cannot play against yourself")
	}
	// TODO: validate Users
	// TODO: validate Clock

	game := domain.Game{
		ID: uuid.NewString(),
		Users: domain.GameUsers{
			White: domain.User{
				ID: e.FormValue("white"),
			},
			Black: domain.User{
				ID: e.FormValue("black"),
			},
		},
		Clock: domain.GameClock{
			Initial:   300,
			Increment: 15,
		},
		Variant: e.FormValue("variant"),
		Winner:  e.FormValue("winner"),
		Status:  e.FormValue("status"),
	}

	// TODO: validate
	log.Default().Printf("creating game: %s", game.ID)
	game, err := a.gameRepo.CreateOrUpdate(ctx, &game)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return item.Game(game).Render(ctx, e.Response().Writer)
}

func (a *api) getGames(e echo.Context) error {
	ctx, cancel := context.WithCancel(e.Request().Context())
	defer cancel()

	games, err := a.gameRepo.GetAll(ctx)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	users, err := a.userRepo.GetAll(ctx)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	c := layout.Base("Games", page.Games(games, users))
	return c.Render(ctx, e.Response().Writer)
}

func (a *api) deleteGame(c echo.Context) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	id := c.Param("id")
	log.Default().Printf("deleting game: %s", id)
	err := a.gameRepo.Delete(ctx, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "deleted game")
}
