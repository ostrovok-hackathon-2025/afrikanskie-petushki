package offer

import (
	"context"
	sql2 "database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

var (
	ErrNotFroundUser = errors.New("user not found")
)

func (o *offer) Create(ctx context.Context, create *model.Create) (uuid.UUID, error) {
	var id uuid.UUID
	sql := `
			INSERT INTO offer (hotel_id, expiration_at, location_id, task)
			VALUES ($1, $2, $3, $4)
			RETURNING id;
			`
	err := o.postgresClient.QueryRowContext(ctx, sql, create.HotelID, create.ExpirationAT, create.LocalID, create.Task).Scan(&id)
	switch {
	case errors.Is(err, sql2.ErrNoRows):
		log.Printf("no user with id %d\n", id)
		return uuid.UUID{}, ErrNotFroundUser
	case err != nil:
		return uuid.UUID{}, err
	default:
		log.Printf("create user with %v\n", id)
	}

	return id, nil
}
