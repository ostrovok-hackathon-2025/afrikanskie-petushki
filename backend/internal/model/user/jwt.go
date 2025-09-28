package user

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID        string `json:"user_id"`
	OstrovokLogin string `json:"ostrovok_login"`
	Email         string `json:"email"`
	IsAdmin       bool   `json:"is_admin"`
	jwt.RegisteredClaims
}
