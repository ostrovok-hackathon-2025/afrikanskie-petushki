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

	GetByOfferIDForDraw(ctx context.Context, offerID uuid.UUID) ([]*application.ApplicationWithRating, error)

	GetUserAppLimitInfo(ctx context.Context, userID uuid.UUID) (*application.UserAppLimitInfo, error)

	UpdateApplicationStatus(ctx context.Context, application *application.Application) error

	GetByFilter(ctx context.Context, filter *application.Filter) ([]*application.Application, error)
	GetCountByFilter(ctx context.Context, filter *application.Filter) (int, error)
}
