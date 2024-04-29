package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/view/layout"
	"github.com/xlund/tracker/internal/view/page"
)

func (a *api) createGameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	Users := domain.GameUsers{
		White: domain.User{
			ID: r.FormValue("white"),
		},
		Black: domain.User{
			ID: r.FormValue("black"),
		},
	}
	if Users.White.ID == Users.Black.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: validate Users

	game := domain.Game{
		ID:    uuid.NewString(),
		Users: Users,
		Clock: domain.GameClock{
			Initial:   300,
			Increment: 15,
		},
		Variant: r.FormValue("variant"),
		Winner:  r.FormValue("winner"),
		Status:  r.FormValue("status"),
	}

	err := a.gameRepo.CreateOrUpdate(ctx, &game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

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
