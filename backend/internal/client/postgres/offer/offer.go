package offer

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

type Offer interface {
	GetByID(ctx context.Context, id string) (*model.Offer, error)
	Get(ctx context.Context, pageSettings *model.PageSettings) ([]*model.Offer, int, error)
	GetByFilter(ctx context.Context, filter *model.Filter) ([]*model.Offer, int, error)

	Create(ctx context.Context, create *model.Create) (uuid.UUID, error)

	Edit(ctx context.Context, filter *model.Edit) error
}

type offer struct {
	postgresClient *sqlx.DB
	logger         *log.Logger
}

func New(postgresClient *sqlx.DB, logger *log.Logger) Offer {
	return &offer{
		postgresClient: postgresClient,
		logger:         logger,
	}
}
