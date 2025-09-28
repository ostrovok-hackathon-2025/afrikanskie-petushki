package user

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func generateTokens(user *model.User, jwtSecret []byte) (*docs.AuthResponse, error) {
	accessToken, err := generateToken(user, 24*time.Hour, "access", jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(user, 30*24*time.Hour, "refresh", jwtSecret)
	if err != nil {
		return nil, err
	}

	return &docs.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessTTL:    86400,   // 24 часа
		RefreshTTL:   2592000, // 30 дней
	}, nil
}

func generateToken(user *model.User, duration time.Duration, audience string, jwtSecret []byte) (string, error) {
	userIDStr := user.ID.String()
	claims := model.JWTClaims{
		UserID:        userIDStr,
		OstrovokLogin: user.OstrovokLogin,
		Email:         user.Email,
		IsAdmin:       user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ostrovok-secret-guest",
			Subject:   userIDStr,
			Audience:  []string{audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
