package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

type RegisterRequest struct {
	OstrovokLogin string `json:"ostrovok_login"`
	Password      string `json:"password"`
}

func (u *useCase) Register(ctx context.Context, req *docs.SignUpRequest) (*docs.AuthResponse, error) {
	// 1. Получаем пользователя из системы Островок
	ostrovokUser, err := u.ostrovokClient.GetUserByLogin(ctx, req.OstrovokLogin)
	if err != nil {
		return nil, err
	}

	exists, err := u.repo.UserExists(ctx, req.OstrovokLogin)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("пользователь уже зарегистрирован")
	}

	passwordHash := pkg.HashPassword(req.Password)

	user := &model.User{
		ID:            uuid.New(),
		OstrovokLogin: ostrovokUser.Login,
		Email:         ostrovokUser.Email,
		IsAdmin:       false,
	}

	err = u.repo.CreateUser(ctx, user, passwordHash)
	if err != nil {
		return nil, err
	}

	return generateTokens(user, u.jwtSecret)
}
