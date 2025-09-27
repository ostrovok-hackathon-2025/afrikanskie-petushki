package register

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
)

type User struct {
	ID            string
	OstrovokLogin string
	Email         string
	IsAdmin       bool
}

type UserRepository interface {
	FindUserByLogin(ctx context.Context, login string) (*User, string, error) // возвращает user, password_hash, error
	CreateUser(ctx context.Context, user *User, passwordHash string) error
	UserExists(ctx context.Context, login string) (bool, error)
}

type JWTClaims struct {
	UserID        string `json:"user_id"`
	OstrovokLogin string `json:"ostrovok_login"`
	Email         string `json:"email"`
	IsAdmin       bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Login(ctx context.Context, req docs.LogInRequest) (*docs.AuthResponse, error) {
	user, storedPasswordHash, err := s.userRepo.FindUserByLogin(ctx, req.OstrovokLogin)
	if err != nil {
		return nil, err
	}

	if hashPassword(req.Password) != storedPasswordHash {
		return nil, errors.New("incorrect password")
	}

	return generateTokens(user, s.jwtSecret)
}

func (s *AuthService) generateToken(user *User, duration time.Duration, audience string) (string, error) {
	return generateToken(user, duration, audience, s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
