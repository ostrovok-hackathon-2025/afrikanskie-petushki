package offer

import (
	"context"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

func (u *useCase) GetByID(ctx context.Context, id uuid.UUID) (model.Offer, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) GetForPage(
	ctx context.Context,
	pageSettings model.PageSettings,
) (offers []model.Offer, pageCount int, err error) {
	filter := model.Filter{
		PageSettings: pageSettings,
		LocationID:   pkg.NewEmpty[uuid.UUID](),
	}
	return u.repo.GetByFilter(ctx, filter)
}

func (u *useCase) GetByFilter(
	ctx context.Context,
	filter model.Filter,
) (offers []model.Offer, pagesCount int, err error) {
	return u.repo.GetByFilter(ctx, filter)
}
