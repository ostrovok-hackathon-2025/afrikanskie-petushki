package ostrovok

import (
	"context"
	"errors"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/ostrovok"
)

var (
	ErrUserNotExists = errors.New("user don't found in ostrovok")
	users            = map[string]*model.OstrovokUser{
		"doverlof": {
			Login: "doverlof",
			Email: "doverlof@exampleemail.com",
		},
		"notblinkyet": {
			Login: "notblinkyet",
			Email: "notblinkyet@exampleemail.com",
		},
		"smokingElk": {
			Login: "smokingElk",
			Email: "chicherin@exampleemail.com",
		},
		"sophistik": {
			Login: "sophistik",
			Email: "chicherin@exampleemail.com",
		},
		"chicherin": {
			Login: "chicherin",
			Email: "chicherin@exampleemail.com",
		},
		"root": {
			Login: "root",
			Email: "root@exampleemail.com",
		},
	}
)

type Client interface {
	GetUserByLogin(ctx context.Context, login string) (*model.OstrovokUser, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) GetUserByLogin(_ context.Context, login string) (*model.OstrovokUser, error) {
	if user, ok := users[login]; ok {
		return user, nil
	}
	return nil, ErrUserNotExists
}
