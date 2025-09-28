package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

type Repo interface {
	FindUserByLogin(ctx context.Context, login string) (*model.User, string, error) // возвращает user, password_hash, error
	CreateUser(ctx context.Context, user *model.User, passwordHash string) error
	UserExists(ctx context.Context, login string) (bool, error)
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
	query := `SELECT id, ostrovok_login, password_hash, is_admin FROM "user" WHERE ostrovok_login = $1`

	var user model.User
	var passwordHash string

	err := r.sqlClient.QueryRowContext(ctx, query, login).Scan(
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
