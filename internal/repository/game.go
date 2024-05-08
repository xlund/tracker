package repository

import (
	"context"
	"log"

	"github.com/xlund/tracker/internal/domain"
)

type postgresGameRepository struct {
	conn Connection
}

func NewPostgresGame(conn Connection) domain.GameRepository {
	return &postgresGameRepository{conn: conn}
}

func (r *postgresGameRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Game, error) {
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gg []domain.Game
	for rows.Next() {
		var g domain.Game
		if err := rows.Scan(&g.ID, &g.Users.ID, &g.Status, &g.Winner, &g.Users.White.ID, &g.Users.Black.ID, &g.Users.White.Username, &g.Users.Black.Username, &g.Users.White.Name, &g.Users.Black.Name); err != nil {
			return nil, err
		}
		gg = append(gg, g)
	}
	return gg, nil
}

func (r *postgresGameRepository) GetById(ctx context.Context, id int) (domain.Game, error) {
	query := `
		SELECT id
		FROM games
		WHERE id = $1`

	gg, err := r.fetch(ctx, query, id)
	if err != nil {
		return domain.Game{}, err
	}
	if len(gg) == 0 {
		return domain.Game{}, domain.ErrNotFound
	}

	return gg[0], nil

}

func (r *postgresGameRepository) GetAll(ctx context.Context) ([]domain.Game, error) {

	query := `
	SELECT
    g.id,
    g.status,
    g.winner,
    COALESCE(wp.username, 'removed') AS white_player_username,
    COALESCE(bp.username, 'removed') AS black_player_username,
	COALESCE(wp.name, 'removed') AS white_player_name,
	COALESCE(bp.name, 'removed') AS black_player_name
	FROM
		games g
	LEFT JOIN
		users wp ON g.white = wp.id
	LEFT JOIN
		users bp ON g.black = bp.id;
		`
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		log.Default().Printf("Game Repository GetAll Error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var gg []domain.Game
	for rows.Next() {
		var g domain.Game
		if err := rows.Scan(&g.ID, &g.Status, &g.Winner, &g.Users.White.Username, &g.Users.Black.Username, &g.Users.White.Name, &g.Users.Black.Name); err != nil {
			log.Default().Printf("Game Repository GetAll Error: %v", err)
			return nil, err
		}
		gg = append(gg, g)
	}
	return gg, nil
}

func (r *postgresGameRepository) CreateOrUpdate(ctx context.Context, g *domain.Game) (domain.Game, error) {
	games_query := `
		INSERT INTO games (id, white, black,  status, variant, winner) VALUES ($1, $2, $3, $4, $5, $6 )
		RETURNING id, white, black,  status, variant, winner;

		`
	err := r.conn.QueryRow(ctx, games_query, g.ID, g.Users.White.ID, g.Users.Black.ID, g.Status, g.Variant, g.Winner).Scan(&g.ID, &g.Users.White.ID, &g.Users.Black.ID, &g.Status, &g.Variant, &g.Winner)

	if err != nil {
		return *g, err
	}

	return *g, err
}

func (r *postgresGameRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM games WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}
