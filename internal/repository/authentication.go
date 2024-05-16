package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/xlund/tracker/internal/domain"
	"golang.org/x/oauth2"
)

type auth0Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func NewAuth0Authenticator() (domain.Authenticator, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://"+os.Getenv("AUTH0_DOMAIN")+"/")

	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}
	return &auth0Authenticator{provider, conf}, nil
}

func (a *auth0Authenticator) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {

	token, err := a.exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *auth0Authenticator) exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := a.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *auth0Authenticator) VerifyIdToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func (a *auth0Authenticator) AuthURL(state string) string {
	return a.AuthCodeURL(state)
}
