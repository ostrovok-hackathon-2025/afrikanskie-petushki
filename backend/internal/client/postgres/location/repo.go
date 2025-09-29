package location

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/location"
)

type Repo interface {
	GetAll(ctx context.Context) ([]model.Location, error)
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

func (r *repo) GetAll(ctx context.Context) ([]model.Location, error) {
	var locations []model.Location
	query, _, err := sq.Select("id", "name").From("location").ToSql()
	if err != nil {
		return nil, err
	}
	err = r.sqlClient.SelectContext(ctx, &locations, query)
	if err != nil {
		//TODO обработка похитрее
		return nil, err
	}
	return locations, nil
}

func (r *repo) Create(ctx context.Context, name string) (uuid.UUID, error) {
	id := uuid.New()
	query, args, err := sq.Insert("location").Columns(
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
