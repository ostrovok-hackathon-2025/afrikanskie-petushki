package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

type ApplicationUseCase interface {
	CreateApplication(ctx context.Context, userId uuid.UUID, offerId uuid.UUID) (uuid.UUID, error)
	GetApplications(ctx context.Context, userId uuid.UUID, pageNum, pageSize int) ([]*application.Application, int, error)
	GetApplicationById(ctx context.Context, userId uuid.UUID, applicationId uuid.UUID) (*application.Application, error)
	GetUserAppLimitInfo(ctx context.Context, userID uuid.UUID) (*application.UserAppLimitInfo, error)
}
