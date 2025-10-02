package hotel

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/hotel"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/hotel"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]model.Hotel, error)
	Create(ctx context.Context, create model.Create) (uuid.UUID, error)
}

type useCase struct {
	repo hotel.Repo
}

func NewUseCase(repo hotel.Repo) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) GetAll(ctx context.Context) ([]model.Hotel, error) {
	return u.repo.GetAll(ctx)
}

func (u *useCase) Create(ctx context.Context, create model.Create) (uuid.UUID, error) {
	return u.repo.Create(ctx, create)
}
