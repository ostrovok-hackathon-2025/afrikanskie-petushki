package offer

import (
	"context"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (u *useCase) GetByID(ctx context.Context, id string) (*model.Offer, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) GetForPage(
	ctx context.Context,
	pageSettings *model.PageSettings,
) (offers []*model.Offer, pageCount int, err error) {
	return u.repo.GetForPage(ctx, pageSettings)
}

func (u *useCase) GetByFilter(
	ctx context.Context,
	filter *model.Filter,
) (offers []*model.Offer, pagesCount int, err error) {
	return u.repo.GetByFilter(ctx, filter)
}
