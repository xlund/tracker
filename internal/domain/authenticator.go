package domain

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Authenticator interface {
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	AuthURL(state string) string
	VerifyIdToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}
