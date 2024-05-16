package api

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/repository"
	"golang.org/x/oauth2"
)

type api struct {
	httpClient *http.Client

	userRepo      domain.UserRepository
	gameRepo      domain.GameRepository
	authenticator domain.Authenticator
}

func NewApi(ctx context.Context, pool *pgxpool.Pool) *api {

	userRepo := repository.NewPostgresUser(pool)

	gameRepo := repository.NewPostgresGame(pool)

	authenticator, _ := repository.NewAuth0Authenticator()

	client := &http.Client{}
	return &api{
		httpClient:    client,
		authenticator: authenticator,

		userRepo: userRepo,
		gameRepo: gameRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *echo.Echo {
	e := echo.New()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})
	gob.Register(&oauth2.Token{})
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("/static", "js")
	e.Static("/static", "css")

	e.GET("/", a.getIndex)

	e.GET("/users", a.getUsers)
	e.POST("/users/new", a.createUser)
	e.DELETE("/users/:id", a.deleteUser)

	e.GET("/users/:id", a.getUser)
	e.GET("/users/me", a.getUser)

	e.GET("/games", a.getGames)
	e.POST("/games/new", a.createGame)
	e.DELETE("/games/:id", a.deleteGame)

	e.GET("/login", a.login)
	e.GET("/authenticated", a.authCallback)
	e.GET("/logout", a.logout)

	e.GET("/v1/health", a.healthCheck)

	return e
}
