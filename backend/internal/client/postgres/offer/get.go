package offer

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
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
		"o.status",
		"o.participants_limit",
		"(SELECT COUNT(*) FROM application a WHERE a.offer_id = o.id) as participants_count",
	).From("offer o").
		Join("hotel h ON o.hotel_id = h.id").
		Join("room r ON o.room_id = r.id")
)

func (r *repo) GetByFilter(
	ctx context.Context,
	filter model.Filter,
) (offers []model.Offer, err error) {

	// SELECT By filter
	sql := baseGetSql
	if id, ok := filter.ID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.id": id})
	}
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	query, args, err := sql.Limit(filter.Limit).Offset(filter.Offset).PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}
	err = r.sqlClient.SelectContext(ctx, &offers, query, args...)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (r *repo) GetCount(ctx context.Context, filter model.Filter) (int, error) {
	var count int
	sql := sq.Select("COUNT(*)").
		From("offer o").
		Join("hotel h ON o.hotel_id = h.id").
		Join("room r ON o.room_id = r.id")
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	query, args, err := sql.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}
	err = r.sqlClient.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repo) GetByExpirationTime(ctx context.Context) (offers []model.Offer, err error) { // поиск истекших оферов
	sql := baseGetSql.Where(sq.LtOrEq{"o.expiration_at": time.Now()}). //истекшие оферы
										Where(sq.Eq{"o.status": "created"}).
										PlaceholderFormat(sq.Dollar).
										Limit(10) // Берем по 10 оферов, позже запихну в воркер. Для равномерной нагрузки. А то можно надорвать мышцу если резко взять большой вес. А в приседе вообще геморрой может выскочить.
	query, args, err := sql.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.sqlClient.SelectContext(ctx, &offers, query, args...)
	if err != nil {
		return nil, err
	}
	// проверили две ошибки на запрос и на подключение
	return offers, nil
}
