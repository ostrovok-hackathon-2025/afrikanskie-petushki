package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

func (u *useCase) Login(ctx context.Context, req *docs.LogInRequest) (*docs.AuthResponse, error) {
	user, storedPasswordHash, err := u.repo.FindUserByLogin(ctx, req.OstrovokLogin)
	if err != nil {
		return nil, err
	}

	if hashPassword(req.Password) != storedPasswordHash {
		return nil, ErrIncorrectPassword
	}

	return generateTokens(user, u.jwtSecret)
}

func (u *useCase) generateToken(user *model.User, duration time.Duration, audience string) (string, error) {
	return generateToken(user, duration, audience, u.jwtSecret)
}

func (u *useCase) ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return u.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
