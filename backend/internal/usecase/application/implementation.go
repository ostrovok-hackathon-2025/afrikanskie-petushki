package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	applicationRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

type ApplicationService struct {
	repo applicationRepo.ApplicationRepo
}

func NewApplicationService(repo applicationRepo.ApplicationRepo) ApplicationUseCase {
	return &ApplicationService{
		repo: repo,
	}
}

func (s *ApplicationService) CreateApplication(
	ctx context.Context,
	userId uuid.UUID,
	offerId uuid.UUID,
) (uuid.UUID, error) {
	newApplication := application.NewApplication(userId, offerId)

	err := s.repo.CreateApplication(ctx, newApplication)

	switch {
	case errors.Is(err, applicationRepo.ErrOfferNotExist) ||
		errors.Is(err, applicationRepo.ErrUserNotExist) ||
		errors.Is(err, applicationRepo.ErrParticipantsLimit) ||
		errors.Is(err, applicationRepo.ErrAppLimit):
		return uuid.UUID{}, err
	case err != nil:
		return uuid.UUID{}, fmt.Errorf("failed to create application in repo: %w", err)
	}

	return newApplication.Id, nil
}

func (s *ApplicationService) GetApplications(
	ctx context.Context,
	userId uuid.UUID,
	pageNum, pageSize int,
) ([]*application.Application, int, error) {
	applications, pagesCount, err := s.repo.GetApplications(ctx, userId, pageNum, pageSize)

	if err == applicationRepo.ErrPageNotFound {
		return nil, 0, err
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all applications from repo: %w", err)
	}

	return applications, pagesCount, nil
}

func (s *ApplicationService) GetApplicationById(
	ctx context.Context,
	userId uuid.UUID,
	applicationId uuid.UUID,
) (*application.Application, error) {
	app, err := s.repo.GetApplicationById(ctx, applicationId)

	if err == applicationRepo.ErrApplicationNotFound {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get application by id from repo: %w", err)
	}

	if app.UserId != userId {
		return nil, ErrNotOwner
	}

	return app, nil
}

func (s *ApplicationService) GetUserAppLimitInfo(ctx context.Context, userID uuid.UUID) (*application.UserAppLimitInfo, error) {
	info, err := s.repo.GetUserAppLimitInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user app limit info from repo: %w", err)
	}
	return info, nil
}
