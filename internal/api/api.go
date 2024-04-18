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
}

func NewApi(ctx context.Context, pool *pgxpool.Pool) *api {

	userRepo := repository.NewPostgresUser(pool)

	client := &http.Client{}
	return &api{
		httpClient: client,

		userRepo: userRepo,
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

	r.HandleFunc("GET /v1/health", a.healthCheckHandler)
	r.HandleFunc("POST /users/new", a.createUserHandler)

	return r
}
