package offer

import (
	"context"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (u *useCase) Edit(ctx context.Context, filter *model.Edit) error {
	return u.repo.Edit(ctx, filter)
}
