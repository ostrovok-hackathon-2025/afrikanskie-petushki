package offer

import (
	"context"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (u *useCase) Create(ctx context.Context, create *model.Create) (uuid.UUID, error) {
	var check model.Check
	_, _ = u.repo.Check(ctx, &check)

	return u.repo.Create(ctx, create)
}
