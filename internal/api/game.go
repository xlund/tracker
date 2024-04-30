package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/view/layout"
	"github.com/xlund/tracker/internal/view/page"
	"github.com/xlund/tracker/internal/view/partial/table/item"
)

func (a *api) createGameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	whiteID := r.FormValue("white")
	blackID := r.FormValue("black")
	if whiteID == blackID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: validate Users
	// TODO: validate Clock

	game := domain.Game{
		ID: uuid.NewString(),
		Users: domain.GameUsers{
			White: domain.User{
				ID: r.FormValue("white"),
			},
			Black: domain.User{
				ID: r.FormValue("black"),
			},
		},
		Clock: domain.GameClock{
			Initial:   300,
			Increment: 15,
		},
		Variant: r.FormValue("variant"),
		Winner:  r.FormValue("winner"),
		Status:  r.FormValue("status"),
	}

	// TODO: validate
	log.Default().Printf("creating game: %s", game.ID)
	game, err := a.gameRepo.CreateOrUpdate(ctx, &game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}
	item.Game(game).Render(ctx, w)

}

func (a *api) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	games, err := a.gameRepo.GetAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

	users, err := a.userRepo.GetAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

	c := layout.Base("Games", page.Games(games, users))
	c.Render(ctx, w)
}

func (a *api) deleteGameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := r.PathValue("id")
	log.Default().Printf("deleting game: %s", id)
	err := a.gameRepo.Delete(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
