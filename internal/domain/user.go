package domain

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID       string
	Username string

	CreatedAt string
	Name      string
}

func (u *User) NormalizedName() string {
	return strings.ToLower(u.Name)
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(3, 32)),
		validation.Field(&u.Name, validation.Required, validation.Length(3, 0)),
	)
}

type UserRepository interface {
	GetById(context.Context, int) (User, error)
	GetByUsername(context.Context, string) (User, error)

	CreateOrUpdate(context.Context, *User) error
	Delete(context.Context, string) error
}
