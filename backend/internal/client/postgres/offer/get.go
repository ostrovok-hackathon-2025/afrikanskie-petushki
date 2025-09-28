package offer

import (
	"context"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (r *repo) GetByID(ctx context.Context, id string) (*model.Offer, error) {
	var res *model.Offer
	sql := `
			SELECT id, hotel_id, location_id, expiration_at, used, task FROM offer WHERE id = $1;
			`
	err := r.sqlClient.GetContext(ctx, res, sql, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) GetForPage(
	ctx context.Context,
	pageSettings *model.PageSettings,
) (offers []*model.Offer, pagesCount int, err error) {
	tr, err := r.sqlClient.Beginx()
	if err != nil {
		return nil, 0, err
	}
	sql := `
			SELECT * FROM offer LIMIT $1, $2;
			`
	err = tr.SelectContext(ctx, offers, sql, pageSettings.Offset, pageSettings.Limit)
	if err != nil {
		//TODO обратка похитрее
		return nil, 0, err
	}
	sql = `
			SELECT COUNT(*) FROM offer;
			`
	var count int
	err = tr.SelectContext(ctx, count, sql)
	if err != nil {
		//TODO обратка похитрее
		return nil, 0, err
	}
	err = tr.Commit()
	if err != nil {
		return nil, 0, err
	}
	if count%pageSettings.Limit == 0 {
		return offers, count / pageSettings.Limit, nil
	}
	return offers, count/pageSettings.Limit + 1, nil
}

func (r *repo) GetByFilter(
	ctx context.Context,
	filter *model.Filter,
) (offers []*model.Offer, pagesCount int, err error) {
	sql := `
			SELECT o.id as offer_id,
			       o.hotel_id as hotel_id,
			       o.room_id as room_id,
			       o.check_in_at as check_in_at,
			       o.check_out_at as check_out_at,
			       o.expiration_at as expiration_at,
			       o.task as task
			FROM offer o JOIN hotel h ON o.hotel_id = h.id WHERE h.location_id=$1 LIMIT $2, $3;
			`
	err = r.sqlClient.SelectContext(ctx, offers, sql, filter.LocationID, filter.PageSettings.Offset, filter.PageSettings.Limit)
	if err != nil {
		//TODO обратка похитрее
		return nil, 0, err
	}
	sql = `
			SELECT COUNT(*) FROM offer;
			`
	var count int
	err = r.sqlClient.SelectContext(ctx, count, sql)
	if err != nil {
		//TODO обратка похитрее
		return nil, 0, err
	}
	if count%filter.PageSettings.Limit == 0 {
		return offers, count / filter.PageSettings.Limit, nil
	}
	return offers, count/filter.PageSettings.Limit + 1, nil
}
