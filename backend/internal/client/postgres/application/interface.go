package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

type ApplicationRepo interface {
	CreateApplication(ctx context.Context, application *application.Application) error
	GetApplications(ctx context.Context, userId uuid.UUID, pageNum, pageSize int) ([]*application.Application, int, error)
	GetApplicationById(ctx context.Context, applicationId uuid.UUID) (*application.Application, error)
	GetByOfferID(ctx context.Context, offerID uuid.UUID) ([]*application.Application, error)
}
