package hotel

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/hotel"
)

type Repo interface {
	GetAll(ctx context.Context) ([]*model.Hotel, error)
	Create(ctx context.Context, create *model.Create) (uuid.UUID, error)
}

type repo struct {
	sqlClient *sqlx.DB
}

func NewRepo(sqlClient *sqlx.DB) Repo {
	return &repo{
		sqlClient: sqlClient,
	}
}

func (r *repo) GetAll(ctx context.Context) ([]*model.Hotel, error) {
	var hotels []*model.Hotel
	sql := `
        SELECT 
    	h.id   AS id,
    	h.name AS name,
    	l.id   AS location_id,
    	l.name AS location_name
		FROM hotel h
		JOIN location l ON h.location_id = l.id;
	`
	err := r.sqlClient.SelectContext(ctx, &hotels, sql)
	if err != nil {
		//TODO обработка похитрее
		return nil, err
	}
	return hotels, nil
}

func (r *repo) Create(ctx context.Context, create *model.Create) (uuid.UUID, error) {
	id := uuid.New()
	sql := "INSERT INTO hotel (id, name, location_id) VALUES ($1, $2, $3)"
	err := r.sqlClient.QueryRowContext(ctx, sql, id, create.Name, create.LocationID).Scan()
	if err != nil {
		//TODO обработка похитрее
		return uuid.UUID{}, err
	}
	return id, nil
}
