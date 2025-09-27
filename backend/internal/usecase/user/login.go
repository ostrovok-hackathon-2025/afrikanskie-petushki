package register

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"os"
	"time"
)

type User struct {
	ID            uuid.UUID
	OstrovokLogin string
	Email         string
	IsAdmin       bool
}

type UserRepository interface {
	FindUserByLogin(ctx context.Context, login string) (*User, string, error) // возвращает user, password_hash, error
	CreateUser(ctx context.Context, user *User, passwordHash string) error
	UserExists(ctx context.Context, login string) (bool, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository() (*PostgresUserRepository, error) {
	host := getEnvWithDefault("POSTGRES_HOST", "localhost")
	port := getEnvWithDefault("POSTGRES_PORT", "5432")
	user := getEnvWithDefault("POSTGRES_USER", "admin")
	password := getEnvWithDefault("POSTGRES_PASSWORD", "admin")
	dbname := getEnvWithDefault("POSTGRES_DB", "ostrovok_secret_guest")
	sslmode := getEnvWithDefault("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &PostgresUserRepository{db: db}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (r *PostgresUserRepository) FindUserByLogin(ctx context.Context, login string) (*User, string, error) {
	query := `SELECT id, ostrovok_login, password_hash, is_admin FROM "user" WHERE ostrovok_login = $1`

	var user User
	var passwordHash string

	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.OstrovokLogin,
		&passwordHash,
		&user.IsAdmin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	return &user, passwordHash, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *User, passwordHash string) error {
	query := `INSERT INTO "user" (id, ostrovok_login, password_hash, is_admin) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.OstrovokLogin, passwordHash, user.IsAdmin)
	return err
}

func (r *PostgresUserRepository) UserExists(ctx context.Context, login string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "user" WHERE ostrovok_login = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, login).Scan(&exists)
	return exists, err
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

func NewAuthService(userRepo UserRepository) *AuthService {
	jwtSecret := getEnvWithDefault("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")

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
		return nil, ErrIncorrectPassword
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

	return nil, ErrInvalidToken
}
