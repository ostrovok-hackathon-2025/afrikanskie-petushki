package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
)

type User struct {
	ID            string `json:"id"`
	OstrovokLogin string `json:"ostrovok_login"`
	Email         string `json:"email"`
	IsAdmin       bool   `json:"is_admin"`
}

type UserRepository interface {
	FindUserByLogin(ctx context.Context, login string) (*User, string, error) // возвращает user, password_hash, error
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

	accessToken, err := s.generateToken(user, 24*time.Hour, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, err
	}

	return &docs.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessTTL:    86400,
		RefreshTTL:   604800,
	}, nil
}

func (s *AuthService) generateToken(user *User, duration time.Duration, audience string) (string, error) {
	claims := JWTClaims{
		UserID:        user.ID,
		OstrovokLogin: user.OstrovokLogin,
		Email:         user.Email,
		IsAdmin:       user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ostrovok-secret-guest",
			Subject:   user.ID,
			Audience:  []string{audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
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

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
