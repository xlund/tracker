package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/xlund/tracker/internal/domain"
)

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := domain.User{
		ID:       uuid.NewString(),
		Username: r.FormValue("username"),
		Name:     r.FormValue("name"),
	}

	println(user.ID)

	err := a.userRepo.CreateOrUpdate(ctx, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}
