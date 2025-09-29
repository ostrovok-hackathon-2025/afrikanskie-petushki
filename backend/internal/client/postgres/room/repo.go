package room

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/room"
)

type Repo interface {
	GetAll(ctx context.Context) ([]model.Room, error)
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

func (r *repo) GetAll(ctx context.Context) (rooms []model.Room, err error) {
	query, _, err := sq.Select("id", "name").From("room").ToSql()
	if err != nil {
		return nil, err
	}
	err = r.sqlClient.SelectContext(ctx, &rooms, query)
	if err != nil {
		//TODO обработка похитрее
		return nil, err
	}
	return rooms, nil
}

func (r *repo) Create(ctx context.Context, name string) (uuid.UUID, error) {
	id := uuid.New()
	query, args, err := sq.Insert("room").Columns(
		"id",
		"name",
	).Values(
		id,
		name,
	).ToSql()
	if err != nil {
		//TODO обработка похитрее
		return uuid.Nil, err
	}
	err = r.sqlClient.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		//TODO обработка похитрее
		return uuid.Nil, err
	}
	return id, nil
}
