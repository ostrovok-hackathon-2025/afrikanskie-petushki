package location

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/location"
)

type Repo interface {
	GetAll(ctx context.Context) ([]*model.Location, error)
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

func (r *repo) GetAll(ctx context.Context) ([]*model.Location, error) {
	var locations []*model.Location
	sql := "SELECT * FROM location"
	err := r.sqlClient.SelectContext(ctx, &locations, sql)
	if err != nil {
		//TODO обработка похитрее
		return nil, err
	}
	return locations, nil
}

func (r *repo) Create(ctx context.Context, name string) (uuid.UUID, error) {
	id := uuid.New()
	sql := "INSERT INTO location (id, name) VALUES ($1, $2)"
	err := r.sqlClient.QueryRowContext(ctx, sql, id, name).Scan()
	if err != nil {
		//TODO обработка похитрее
		return uuid.UUID{}, err
	}
	return id, nil
}
