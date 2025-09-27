package register

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
)

// OstrovokUser представляет пользователя из системы Островок
type OstrovokUser struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

// OstrovokService интерфейс для работы с API Островок
type OstrovokService interface {
	GetUserByLogin(ctx context.Context, login string) (*OstrovokUser, error)
}

type RegistrationService struct {
	userRepo    UserRepository
	ostrovokSvc OstrovokService
	jwtSecret   []byte
}

type RegisterRequest struct {
	OstrovokLogin string `json:"ostrovok_login"`
	Password      string `json:"password"`
}

func NewRegistrationService(userRepo UserRepository, ostrovokSvc OstrovokService) *RegistrationService {
	jwtSecret := getEnvWithDefault("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")

	return &RegistrationService{
		userRepo:    userRepo,
		ostrovokSvc: ostrovokSvc,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (s *RegistrationService) Register(ctx context.Context, req RegisterRequest) (*docs.AuthResponse, error) {
	// 1. Получаем пользователя из системы Островок
	ostrovokUser, err := s.ostrovokSvc.GetUserByLogin(ctx, req.OstrovokLogin)
	if err != nil {
		return nil, err
	}

	exists, err := s.userRepo.UserExists(ctx, req.OstrovokLogin)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("пользователь уже зарегистрирован")
	}

	passwordHash := hashPassword(req.Password)

	user := &User{
		ID:            uuid.New(),
		OstrovokLogin: ostrovokUser.Login,
		Email:         ostrovokUser.Email,
		IsAdmin:       false,
	}

	err = s.userRepo.CreateUser(ctx, user, passwordHash)
	if err != nil {
		return nil, err
	}

	return generateTokens(user, s.jwtSecret)
}
