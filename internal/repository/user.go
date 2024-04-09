package repository

import (
	"context"

	"github.com/xlund/tracker/internal/domain"
)

type postgresUserRepository struct {
	conn Connection
}

func NewPostgresUser(conn Connection) domain.UserRepository {
	return &postgresUserRepository{conn: conn}
}

func (r *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uu []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Name); err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}
	return uu, nil
}

func (r *postgresUserRepository) GetById(ctx context.Context, id int) (domain.User, error) {
	query := `
		SELECT id, username, name
		FROM users
		WHERE id = $1`

	uu, err := r.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}
	if len(uu) == 0 {
		return domain.User{}, domain.ErrNotFound
	}

	return uu[0], nil

}

func (r *postgresUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `
		SELECT id, username, name
		FROM users
		WHERE username = $1`

	uu, err := r.fetch(ctx, query, username)
	if err != nil {
		return domain.User{}, err
	}
	if len(uu) == 0 {
		return domain.User{}, domain.ErrNotFound
	}

	return uu[0], nil
}

func (r *postgresUserRepository) CreateOrUpdate(ctx context.Context, u *domain.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO users (username, name, created_at)
		VALUES ($1, $2, Now())
		ON CONFLICT (username) DO NOTHING
		RETURNING id`
	return r.conn.QueryRow(ctx, query, u.Username, u.Name).Scan(&u.ID)
}

func (r *postgresUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}
