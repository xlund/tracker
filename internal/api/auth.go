package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (a *api) login(c echo.Context) error {
	state, err := generateRandomState()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// save state to session
	session, err := session.Get("auth-session", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	session.Values["state"] = state
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	url := a.authenticator.AuthURL(state)
	return c.Redirect(http.StatusFound, url)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)
	return state, nil
}

func (a *api) authCallback(c echo.Context) error {

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	session, err := session.Get("auth-session", c)
	if c.QueryParam("state") != session.Values["state"] {
		return c.String(http.StatusBadRequest, "Invalid state parameter.")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	token, err := a.authenticator.GetToken(ctx, c.QueryParam("code"))
	log.Default().Println(token)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	idToken, err := a.authenticator.VerifyIdToken(ctx, token)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var profile map[string]interface{}

	err = idToken.Claims(&profile)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	session.Values["profile"] = profile
	session.Values["access_token"] = token.AccessToken

	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.Redirect(http.StatusFound, "/users/me")
}

func (a *api) logout(c echo.Context) error {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	returnTo, err := url.Parse(scheme + "://" + c.Request().Host)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	params := url.Values{}
	params.Add("returnTo", returnTo.String())
	params.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = params.Encode()
	return c.Redirect(http.StatusFound, logoutUrl.String())
}
