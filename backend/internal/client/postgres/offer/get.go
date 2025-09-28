package offer

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

var (
	baseGetSql = sq.Select(
		"o.id as offer_id",
		"o.hotel_id as hotel_id",
		"h.name as hotel_name",
		"o.room_id as room_id",
		"r.name as room_name",
		"o.check_in_at",
		"o.check_out_at",
		"o.expiration_at",
		"o.task",
	).From("offer o").
		Join("hotel h ON o.hotel_id = h.id").
		Join("room r ON o.room_id = r.id")
)

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (model.Offer, error) {
	var offer model.Offer
	query, args, err := baseGetSql.Where(sq.Eq{"o.id": id}).ToSql()
	if err != nil {
		return model.Offer{}, err
	}
	err = r.sqlClient.GetContext(ctx, &offer, query, args...)
	if err != nil {
		return model.Offer{}, err
	}
	return offer, nil
}

func (r *repo) GetByFilter(
	ctx context.Context,
	filter model.Filter,
) (offers []model.Offer, pagesCount int, err error) {

	// SELECT By filter
	sql := baseGetSql
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	query, args, err := sql.Limit(uint64(filter.PageSettings.Limit)).Offset(uint64(filter.PageSettings.Offset)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = r.sqlClient.SelectContext(ctx, &offers, query, args...)
	if err != nil {
		return nil, 0, err
	}

	// GET count
	var count int

	sql = sq.Select("COUNT(*)")
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	query, args, err = sql.ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = r.sqlClient.GetContext(ctx, &count, query, args...)
	err = r.sqlClient.SelectContext(ctx, &offers, query, args...)
	if err != nil {
		return nil, 0, err
	}
	if count%filter.PageSettings.Limit == 0 {
		return offers, count / filter.PageSettings.Limit, nil
	}
	return offers, count/filter.PageSettings.Limit + 1, nil
}
