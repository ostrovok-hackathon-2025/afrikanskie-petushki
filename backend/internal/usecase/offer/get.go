package offer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
	"github.com/pkg/errors"
)

var (
	ErrWrongNumOffers = errors.New("get wrong amount of offers found")
)

func (u *useCase) GetByID(ctx context.Context, id uuid.UUID) (model.Offer, error) {
	filter := model.Filter{ID: pkg.NewWithValue(id)}
	offers, err := u.repo.GetByFilter(ctx, filter)
	if err != nil {
		return model.Offer{}, err
	}
	if len(offers) != 1 {
		return model.Offer{}, errors.Wrap(ErrWrongNumOffers, fmt.Sprintf("got %d offers, want: 1", len(offers)))
	}
	return offers[0], nil
}

func (u *useCase) GetForPage(
	ctx context.Context,
	pageSettings model.PageSettings,
) (offers []model.Offer, pageCount int, err error) {

	//GET offers by filter use only limit and offset
	filter := model.Filter{
		Limit:  pkg.NewWithValue(pageSettings.Limit),
		Offset: pkg.NewWithValue(pageSettings.Offset),
	}
	offers, err = u.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	//GET count
	count, err := u.repo.GetCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return offers, count, nil
}

func (u *useCase) GetByFilter(
	ctx context.Context,
	filter model.Filter,
) (offers []model.Offer, pagesCount int, err error) {
	offers, err = u.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	//GET count
	count, err := u.repo.GetCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return offers, count, nil
}
