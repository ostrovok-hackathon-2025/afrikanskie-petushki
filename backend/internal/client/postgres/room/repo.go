package room

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/room"
)

type Repo interface {
	GetAll(ctx context.Context) ([]*model.Room, error)
	Create(ctx context.Context, name string) (uuid.UUID, error)
}

type repo struct {
	sqlClient *sqlx.DB
}

func NewRepo(sqlClient *sqlx.DB) Repo {
	return &repo{
		sqlClient: sqlClient,
	}
}

func (r *repo) GetAll(ctx context.Context) ([]*model.Room, error) {
	var rooms []*model.Room
	sql := "SELECT * FROM room"
	err := r.sqlClient.SelectContext(ctx, &rooms, sql)
	if err != nil {
		//TODO обработка похитрее
		return nil, err
	}
	return rooms, nil
}

func (r *repo) Create(ctx context.Context, name string) (uuid.UUID, error) {
	id := uuid.New()
	sql := "INSERT INTO room (id, name) VALUES ($1, $2)"
	err := r.sqlClient.QueryRowContext(ctx, sql, id, name).Scan()
	if err != nil {
		//TODO обработка похитрее
		return uuid.UUID{}, err
	}
	return id, nil
}
