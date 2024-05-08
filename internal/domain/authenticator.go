package domain

import (
	"context"
	"net/http"
)

type Authenticator interface {
	// Not sure if it needs to return a user
	CreateUser(ctx context.Context, u User) (string, error)
	StartWebAuthnRegister(ctx context.Context, u User, r *http.Request) (string, error)
	CompleteWebAuthRegister(ctx context.Context, r *http.Request) error
	RemoveUser(ctx context.Context, id string) error
}
