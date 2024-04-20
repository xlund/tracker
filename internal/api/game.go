package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/xlund/tracker/internal/domain"
)

func (a *api) createGameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	game := domain.Game{
		ID: uuid.NewString(),
		Players: domain.GamePlayers{
			White: domain.User{
				ID: "20f669a4-1600-41c5-943e-e1bbdfc0d0fa"},

			Black: domain.User{
				ID: "371703ab-597b-42c0-8199-88f0512206a3"},
		},
		Clock: domain.GameClock{
			Initial:   300,
			Increment: 15,
		},
		Source:  r.FormValue("source"),
		Variant: r.FormValue("variant"),
		Winner:  r.FormValue("winner"),
		Status:  r.FormValue("status"),
	}

	err := a.gameRepo.CreateOrUpdate(ctx, &game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)

}
