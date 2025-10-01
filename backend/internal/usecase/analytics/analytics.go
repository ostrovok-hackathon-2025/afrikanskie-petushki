package analytics

import (
	"context"

	analyticsRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/analytics"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/analytics"
)

type UseCase interface {
	GetAnalytics(ctx context.Context) (analytics.Analytics, error)
}

type useCase struct {
	repo analyticsRepo.Repo
}

func NewAnalyticsUseCase(repo analyticsRepo.Repo) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) GetAnalytics(ctx context.Context) (analytics.Analytics, error) {
	return u.repo.GetAnalytics(ctx)
}
