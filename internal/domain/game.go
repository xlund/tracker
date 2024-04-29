package domain

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrGameNotFound = fmt.Errorf("game not found")
)

type Game struct {
	ID        string
	Users     GameUsers
	Clock     GameClock
	Source    string
	Status    string
	Variant   string
	Winner    string
	CreatedAt string
}

func (g *Game) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Users, validation.Required),
		validation.Field(&g.Clock, validation.Required),
		validation.Field(&g.Source, validation.Required),
		validation.Field(&g.Status, validation.Required),
		validation.Field(&g.Variant, validation.Required),
	)
}

type GameUsers struct {
	ID    int
	White User
	Black User
}

type GameClock struct {
	Initial   int
	Increment int
}

func (c *GameClock) Readable() string {
	var s string
	s = fmt.Sprintf("%d|%d", c.Initial, c.Increment)
	if c.Initial > 60 {
		s = fmt.Sprintf("%d:%02d|%d", c.Initial/60, c.Initial%60, c.Increment)
	}
	if c.Initial > 60*60 {
		s = fmt.Sprintf("%d:%02d:%02d|%d", c.Initial/60/60, c.Initial/60%60, c.Initial%60, c.Increment)
	}
	return s
}

type GameRepository interface {
	GetById(context.Context, int) (Game, error)
	GetAll(context.Context) ([]Game, error)

	CreateOrUpdate(context.Context, *Game) error
	Delete(context.Context, string) error
}
