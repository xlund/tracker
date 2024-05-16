package domain

import (
	"context"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID       string
	Username string
	Email    string

	CreatedAt string
	AuthID    string
	Name      string
}

type UserWithGames struct {
	User
	Games []Game
}

type UserAuthProfile struct {
	Nickname  string
	Picture   string
	Name      string
	GivenName string
}

func (u *User) NormalizedName() string {
	return strings.ToLower(u.Name)
}

func UserAUthProfileFromMap(m map[string]interface{}) UserAuthProfile {

	return UserAuthProfile{
		Nickname:  m["nickname"].(string),
		Picture:   m["picture"].(string),
		Name:      m["name"].(string),
		GivenName: m["given_name"].(string),
	}
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(3, 32)),
		validation.Field(&u.Name, validation.Length(3, 0)),
		validation.Field(&u.Email, validation.Required, validation.Length(3, 0),
			validation.Match(regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`))),
	)
}

type UserRepository interface {
	GetAll(context.Context) ([]User, error)
	GetById(context.Context, string) (User, error)
	GetByUsername(context.Context, string) (User, error)

	GetByIdWithGames(context.Context, string) (UserWithGames, error)

	Search(context.Context, string) ([]User, error)

	CreateOrUpdate(context.Context, *User) error
	Delete(context.Context, string) error
}
