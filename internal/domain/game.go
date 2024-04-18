package domain

import (
	"context"
	"fmt"
	"time"
)

var (
	ErrGameNotFound = fmt.Errorf("game not found")
)

type Game struct {
	ID          int
	WhitePlayer User
	BlackPlayer User
	TimeControl TimeControl
	Winner      string
	Result      string
	CreatedAt   string
	Name        string
}

type TimeControl struct {
	Start     time.Duration
	Increment time.Duration
}

func (t *TimeControl) Readable() string {
	return fmt.Sprintf("%d|%d", t.Start, t.Increment)
}

type GameRepository interface {
	GetById(context.Context, int) (User, error)
	GetByUsername(context.Context, string) (User, error)

	CreateOrUpdate(context.Context, *User) error
	Delete(context.Context, int) error
}
