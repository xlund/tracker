package cmdutil

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabasePool(ctx context.Context, maxConns int) (*pgxpool.Pool, error) {
	if maxConns == 0 {
		maxConns = 1
	}
	url := fmt.Sprintf("%s?pool_max_conns=%d&pool_min_conns=%d", os.Getenv("DATABASE_URL"), maxConns, 2)

	config, err := pgxpool.ParseConfig(url)

	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, config)
}
