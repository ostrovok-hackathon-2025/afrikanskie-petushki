package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/ostrovok"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/achievement"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

type UseCase interface {
	Register(ctx context.Context, req *docs.SignUpRequest) (*docs.AuthResponse, error)
	ValidateToken(tokenString string) (*model.JWTClaims, error)
	Login(ctx context.Context, req *docs.LogInRequest) (*docs.AuthResponse, error)
	GetMe(ctx context.Context, userId uuid.UUID) (*model.User, error)
}

type useCase struct {
	repo            user.Repo
	ostrovokClient  ostrovok.Client
	jwtSecret       []byte
	achievementRepo achievement.Repo
}

func NewUseCase(repo user.Repo, ostrovokClient ostrovok.Client, achievementRepo achievement.Repo) UseCase {

	jwtSecret := getEnvWithDefault("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")

	return &useCase{
		repo:            repo,
		ostrovokClient:  ostrovokClient,
		jwtSecret:       []byte(jwtSecret),
		achievementRepo: achievementRepo,
	}
}
