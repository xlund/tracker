package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xlund/tracker/internal/domain"
	"github.com/xlund/tracker/internal/repository"
)

type api struct {
	httpClient *http.Client

	userRepo domain.UserRepository
	gameRepo domain.GameRepository
}

func NewApi(ctx context.Context, pool *pgxpool.Pool) *api {

	userRepo := repository.NewPostgresUser(pool)

	gameRepo := repository.NewPostgresGame(pool)

	client := &http.Client{}
	return &api{
		httpClient: client,

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

func (a *api) Routes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("../static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	r.HandleFunc("/", a.getIndexHandler)
	r.HandleFunc("GET /users", a.getUsersHandler)
	r.HandleFunc("POST /users/new", a.createUserHandler)
	r.HandleFunc("POST /users/search", a.searchUsersHandler)
	r.HandleFunc("GET /users/{id}", a.getUserHandler)
	r.HandleFunc("DELETE /users/{id}", a.deleteUserHandler)
	r.HandleFunc("POST /games/new", a.createGameHandler)
	r.HandleFunc("GET /games", a.getGamesHandler)
	r.HandleFunc("DELETE /games/{id}", a.deleteGameHandler)
	r.HandleFunc("GET /v1/health", a.healthCheckHandler)

	return r
}
