package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

const pgForeginKeyErr = "foreign_key_violation"

type applicationRepo struct {
	db *sqlx.DB
}

func NewApplicationRepo(db *sqlx.DB) ApplicationRepo {
	return &applicationRepo{
		db: db,
	}
}

func (r *applicationRepo) CreateApplication(
	ctx context.Context,
	application *application.Application,
) error {
	query := `
	INSERT INTO application (id, user_id, offer_id, status) VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, application.Id, application.UserId, application.OfferId, application.Status)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == pgForeginKeyErr {
			return ErrOfferNotExist
		}

		return fmt.Errorf("failed to insert application into db: %w", err)
	}

	return nil
}

func (r *applicationRepo) GetApplications(
	ctx context.Context,
	userId uuid.UUID,
	pageNum, pageSize int,
) ([]*application.Application, int, error) {
	offset := pageNum * pageSize

	query := `
	SELECT a.id, a.user_id, a.offer_id, a.status, o.expiration_at, h.name FROM application as a
	INNER JOIN offer as o ON a.offer_id = o.id
	INNER JOIN hotel as h ON o.hotel_id = h.id
	WHERE a.user_id = $1
	OFFSET $2
	LIMIT $3
	`

	var apps []ApplicationDTO

	err := r.db.SelectContext(ctx, &apps, query, userId, offset, pageSize)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if pageNum == 0 {
				return []*application.Application{}, 0, nil
			}

			return nil, 0, ErrPageNotFound
		}

		return nil, 0, fmt.Errorf("failed to get all applications from db: %w", err)
	}

	var count int
	query = "SELECT COUNT(*) FROM application WHERE user_id = $1"

	err = r.db.Get(&count, query, userId)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count of applications: %w", err)
	}

	pagesCount := count / pageSize
	if count%pageSize != 0 {
		pagesCount++
	}

	res := make([]*application.Application, 0, len(apps))

	for _, app := range apps {
		res = append(res, app.ToApplicationModel())
	}

	return res, pagesCount, nil
}

func (r *applicationRepo) GetApplicationById(
	ctx context.Context,
	applicationId uuid.UUID,
) (*application.Application, error) {
	query := `
	SELECT a.id, a.user_id, a.offer_id, a.status, o.expiration_at, h.name FROM application as a
	INNER JOIN offer as o ON a.offer_id = o.id
	INNER JOIN hotel as h ON o.hotel_id = h.id
	WHERE a.id = $1
	`

	var app ApplicationDTO

	err := r.db.GetContext(ctx, &app, query, applicationId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}

		return nil, fmt.Errorf("failed to get application from repo: %w", err)
	}

	res := app.ToApplicationModel()

	return res, nil
}
