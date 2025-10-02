package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/validation"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

type Repo interface {
	FindUserByLogin(ctx context.Context, login string) (*model.User, string, error) // возвращает user, password_hash, error
	CreateUser(ctx context.Context, user *model.User, passwordHash string) error
	UserExists(ctx context.Context, login string) (bool, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error)
	GetUserByReportId(ctx context.Context, reportId uuid.UUID) (*model.User, error)
	UpdateRating(ctx context.Context, userId uuid.UUID, rating int) error
}

type repo struct {
	sqlClient *sqlx.DB
}

func NewUserRepo(sqlClient *sqlx.DB) Repo {
	return &repo{
		sqlClient: sqlClient,
	}
}

func (r *repo) FindUserByLogin(ctx context.Context, login string) (*model.User, string, error) {
	query := `SELECT id, ostrovok_login, password_hash, is_admin, rating FROM "user" WHERE ostrovok_login = $1`

	var user model.User
	var passwordHash string

	err := r.sqlClient.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.OstrovokLogin,
		&passwordHash,
		&user.IsAdmin,
		&user.Rating,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	return &user, passwordHash, nil
}

func (r *repo) CreateUser(ctx context.Context, user *model.User, passwordHash string) error {
	query := `INSERT INTO "user" (id, ostrovok_login, password_hash, is_admin) VALUES ($1, $2, $3, $4)`

	_, err := r.sqlClient.ExecContext(ctx, query, user.ID, user.OstrovokLogin, passwordHash, user.IsAdmin)
	return err
}

func (r *repo) UserExists(ctx context.Context, login string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "user" WHERE ostrovok_login = $1)`

	var exists bool
	err := r.sqlClient.QueryRowContext(ctx, query, login).Scan(&exists)
	return exists, err
}

func (r *repo) GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	query := `SELECT id, ostrovok_login, is_admin, rating FROM "user" WHERE id = $1`

	var user UserDTO

	err := r.sqlClient.GetContext(ctx, &user, query, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user from repo: %w", err)
	}

	return user.ToUserModel(), nil
}

func (r *repo) GetUserByReportId(ctx context.Context, reportId uuid.UUID) (*model.User, error) {
	query := `
		SELECT u.id, u.ostrovok_login, u.is_admin, u.rating 
		FROM "user" u 
		INNER JOIN application a ON a.user_id = u.id 
		INNER JOIN report r ON r.application_id = a.id
		WHERE r.id = $1
	`

	var user UserDTO

	err := r.sqlClient.GetContext(ctx, &user, query, reportId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user from repo: %w", err)
	}

	return user.ToUserModel(), nil
}

func (r *repo) UpdateRating(ctx context.Context, userId uuid.UUID, rating int) error {

	rating = validation.ValidateRating(rating)

	query := `UPDATE "user" SET rating = $1 WHERE id = $2`

	result, err := r.sqlClient.ExecContext(ctx, query, rating, userId)
	if err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
