package offer

import (
	"context"
	sql2 "database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

var (
	ErrNotFroundUser = errors.New("user not found")
)

func (r *repo) Create(ctx context.Context, id uuid.UUID, create model.Create) error {
	query, args, err := sq.Insert("offer").Columns(
		"id",
		"hotel_id",
		"room_id",
		"check_in_at",
		"check_out_at",
		"expiration_at",
		"task",
		"participants_limit",
	).Values(
		id,
		create.HotelID,
		create.RoomID,
		create.CheckIn,
		create.CheckOut,
		create.ExpirationAT,
		create.Task,
		create.ParticipantsLimit,
	).PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}
	_, err = r.sqlClient.ExecContext(ctx, query, args...)
	switch {
	case errors.Is(err, sql2.ErrNoRows):
		log.Printf("no user with id %d\n", id)
		return ErrNotFroundUser
	case err != nil:
		return err
	default:
		log.Printf("create user with %v\n", id)
	}
	return nil
}
