package offer

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

type Repo interface {
	GetByID(ctx context.Context, id string) (*model.Offer, error)
	GetForPage(ctx context.Context, pageSettings *model.PageSettings) ([]*model.Offer, int, error)
	GetByFilter(ctx context.Context, filter *model.Filter) ([]*model.Offer, int, error)

	Check(ctx context.Context, check *model.Check) (bool, error)

	Create(ctx context.Context, create *model.Create) (uuid.UUID, error)

	Edit(ctx context.Context, filter *model.Edit) error
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
