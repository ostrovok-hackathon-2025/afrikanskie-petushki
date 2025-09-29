package offer

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

type Repo interface {
	GetByFilter(ctx context.Context, filter model.Filter) ([]model.Offer, error)
	GetCount(ctx context.Context, filter model.Filter) (int, error)

	Create(ctx context.Context, id uuid.UUID, create model.Create) error

	Edit(ctx context.Context, filter model.Edit) error
}

type repo struct {
	sqlClient *sqlx.DB
	logger    *log.Logger
}

func New(sqlClient *sqlx.DB, logger *log.Logger) Repo {
	return &repo{
		sqlClient: sqlClient,
		logger:    logger,
	}
}
