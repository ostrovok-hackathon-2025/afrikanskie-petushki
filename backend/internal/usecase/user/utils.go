package register

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"time"
)

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func generateUserID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func generateTokens(user *User, jwtSecret []byte) (*docs.AuthResponse, error) {
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

func generateToken(user *User, duration time.Duration, audience string, jwtSecret []byte) (string, error) {
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
	return token.SignedString(jwtSecret)
}
