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

func (r *repo) Create(ctx context.Context, create model.Create) (uuid.UUID, error) {
	id := uuid.New()
	sql := `
			INSERT INTO offer (id, hotel_id, room_id, check_in_at, check_out_at, expiration_at, task)
			VALUES ($1, $2, $3, $4, $5, $6, $7);
			`
	err := r.sqlClient.QueryRowContext(
		ctx, sql, id, create.HotelID,
		create.RoomID, create.CheckIn,
		create.CheckOut, create.ExpirationAT, create.Task,
	).Scan()
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
