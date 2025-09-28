package offer

import (
	"context"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (u *useCase) Create(ctx context.Context, create model.Create) (uuid.UUID, error) {
	id := uuid.New()
	err := u.repo.Create(ctx, id, create)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
