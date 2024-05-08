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

func (a *api) authHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	c := layout.Base("Login", page.Auth())

	c.Render(ctx, w)
}

func (a *api) beginRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	log.Default().Println("beginRegistrationHandler")
	username := r.FormValue("username")
	email := r.FormValue("email")

	user := domain.User{
		ID:       uuid.NewString(),
		Username: username,
		Email:    email,
	}

	log.Default().Printf("beginRegistrationHandler::user: %+v", user)

	challenge, err := a.authenticator.StartWebAuthnRegister(ctx, user, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

	_, err = w.Write([]byte(challenge))
	if err != nil {
		log.Default().Println(err.Error())
		return
	}

}

func (a *api) completeRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	log.Default().Println("completeRegistrationHandler")

	err := a.authenticator.CompleteWebAuthRegister(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

	log.Default().Println("registration complete")

	w.WriteHeader(http.StatusOK)
}
