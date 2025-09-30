package offer

import (
	"context"

	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (r *repo) Edit(ctx context.Context, edit model.Edit) error {
	sql := sq.Update("offer")
	if task, ok := edit.Task.Get(); ok {
		sql = sql.Set("task", task)
	}
	if roomID, ok := edit.RoomID.Get(); ok {
		sql = sql.Set("room_id", roomID)
	}
	if hotelID, ok := edit.HotelID.Get(); ok {
		sql = sql.Set("hotel_id", hotelID)
	}
	if checkIn, ok := edit.CheckIn.Get(); ok {
		sql = sql.Set("check_in", checkIn)
	}
	if checkOut, ok := edit.CheckOut.Get(); ok {
		sql = sql.Set("check_out", checkOut)
	}
	if expirationAt, ok := edit.ExpirationAT.Get(); ok {
		sql = sql.Set("expiration_at", expirationAt)
	}
	query, args, err := sql.Where(sq.Eq{"id": edit.OfferID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	err = r.sqlClient.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) EditStatus(ctx context.Context, offerID uuid.UUID, status string) error {
	query, args, err := sq.Update("offer").
		Set("status", status).
		Where(sq.Eq{"id": offerID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.sqlClient.ExecContext(ctx, query, args...)
	return err
}
