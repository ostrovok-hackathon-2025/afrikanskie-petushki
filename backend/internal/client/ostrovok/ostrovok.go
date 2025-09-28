package ostrovok

import (
	"context"
	"errors"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/ostrovok"
)

var (
	ErrUserNotExists = errors.New("user don't found in ostrovok")
	users            = map[string]*model.OstrovokUser{
		"chicherin": {
			Login: "chicherin",
			Email: "yarik_vodila@gmail.com",
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

func (c *client) GetUserByLogin(ctx context.Context, login string) (*model.OstrovokUser, error) {
	if user, ok := users[login]; ok {
		return user, nil
	}
	return nil, ErrUserNotExists
}
