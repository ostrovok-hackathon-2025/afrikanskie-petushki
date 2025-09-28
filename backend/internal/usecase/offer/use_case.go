package offer

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

type UseCase interface {
	Create(ctx context.Context, create *model.Create) (uuid.UUID, error)

	GetByID(ctx context.Context, id string) (*model.Offer, error)
	GetForPage(ctx context.Context, pageSettings *model.PageSettings) (offers []*model.Offer, pageCount int, err error)
	GetByFilter(ctx context.Context, filter *model.Filter) (offers []*model.Offer, pagesCount int, err error)

	Edit(ctx context.Context, filter *model.Edit) error
}

type useCase struct {
	repo offer.Repo
}

func NewUseCase(repo offer.Repo) UseCase {
	return &useCase{
		repo: repo,
	}
}
