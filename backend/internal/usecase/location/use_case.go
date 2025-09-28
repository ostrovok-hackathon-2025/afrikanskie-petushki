package location

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/location"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/location"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]*model.Location, error)
	Create(ctx context.Context, name string) (uuid.UUID, error)
}

type useCase struct {
	repo location.Repo
}

func NewUseCase(repo location.Repo) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) GetAll(ctx context.Context) ([]*model.Location, error) {
	return u.repo.GetAll(ctx)
}

func (u *useCase) Create(ctx context.Context, name string) (uuid.UUID, error) {
	return u.repo.Create(ctx, name)
}
