package room

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/room"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/room"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]model.Room, error)
	Create(ctx context.Context, name string) (uuid.UUID, error)
}

type useCase struct {
	repo room.Repo
}

func NewUseCase(repo room.Repo) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) GetAll(ctx context.Context) ([]model.Room, error) {
	return u.repo.GetAll(ctx)
}

func (u *useCase) Create(ctx context.Context, name string) (uuid.UUID, error) {
	return u.repo.Create(ctx, name)
}
