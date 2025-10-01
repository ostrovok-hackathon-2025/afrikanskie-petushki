package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/validation"

	"github.com/google/uuid"
	repo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

func (u *useCase) GetMe(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	user, err := u.repo.GetUserById(ctx, userId)

	if errors.Is(err, repo.ErrUserNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user from repo: %w", err)
	}

	validation.ValidateRating(user.Rating)

	return user, err
}
