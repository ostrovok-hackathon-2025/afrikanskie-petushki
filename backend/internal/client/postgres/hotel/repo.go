package hotel

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/hotel"
)

type Repo interface {
	GetAll(ctx context.Context) ([]model.Hotel, error)
	Create(ctx context.Context, create model.Create) (uuid.UUID, error)
}

type repo struct {
	sqlClient *sqlx.DB
}

func NewRepo(sqlClient *sqlx.DB) Repo {
	return &repo{
		sqlClient: sqlClient,
	}
}

func (r *repo) GetAll(ctx context.Context) (hotels []model.Hotel, err error) {
	query, _, err := sq.Select(
		"h.id AS id",
		"h.name AS name",
		"l.id AS location_id",
		"l.name AS location_name",
	).From("hotel h").Join("location l ON h.location_id = l.id").ToSql()
	if err != nil {
		return nil, err
	}
	err = r.sqlClient.SelectContext(ctx, &hotels, query)
	if err != nil {
		//TODO обработка похитрее
		return nil, err
	}
	return hotels, nil
}

func (r *repo) Create(ctx context.Context, create model.Create) (uuid.UUID, error) {
	id := uuid.New()
	query, args, err := sq.Insert("hotel").Columns(
		"id",
		"name",
		"location_id",
	).Values(
		id,
		create.Name,
		create.LocationID,
	).ToSql()
	if err != nil {
		return uuid.Nil, err
	}
	err = r.sqlClient.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		//TODO обработка похитрее
		return uuid.Nil, err
	}
	return id, nil
}
