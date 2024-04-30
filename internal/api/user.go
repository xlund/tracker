package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/view/layout"
	"github.com/xlund/tracker/internal/view/page"
	view "github.com/xlund/tracker/internal/view/page"
	"github.com/xlund/tracker/internal/view/partial/form"
)

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := domain.User{
		ID:       uuid.NewString(),
		Username: r.FormValue("username"),
		Name:     r.FormValue("name"),
	}

	err := a.userRepo.CreateOrUpdate(ctx, &user)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

}

func (a *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	users, err := a.userRepo.GetAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

	c := layout.Base("Users", view.Users(users))
	c.Render(ctx, w)
}

func (a *api) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := r.PathValue("id")
	log.Default().Printf("deleting user %s", r.FormValue("id"))
	err := a.userRepo.Delete(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) searchUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	users, err := a.userRepo.Search(ctx, r.FormValue("q"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}
	c := layout.Base("Users", form.SelectUser(r.FormValue("target"), users))
	c.Render(ctx, w)
}

func (a *api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	id := r.PathValue("id")
	userWithGames, err := a.userRepo.GetByIdWithGames(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}
	c := layout.Base("User", page.UserWithGames(userWithGames))
	c.Render(ctx, w)
}
