package offer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

var (
	ErrAlreadyExists = errors.New("offer: already exists")
)

func (u *useCase) Create(ctx context.Context, create model.Create) (uuid.UUID, error) {
	return u.repo.Create(ctx, create)
}
