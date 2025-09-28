package offer

import (
	"context"

	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

func (r *repo) Edit(_ context.Context, _ *model.Edit) error {
	panic("implement me")
}

//TODO create sql from squirrel
//func (r *repo) Edit(ctx context.Context, filter *model.Edit) error {
//	sql := `
//			UPDATE offer
//			SET expiration_at=$1, task=$2
//			WHERE id=$2;
//			`
//	err := r.sqlClient.QueryRowContext(ctx, sql, filter.ExpirationAT, filter.Task, filter.OfferID).Scan()
//	switch {
//	case errors.Is(err, sql2.ErrNoRows):
//		log.Printf("no user with id %d\n", filter.OfferID)
//		return ErrNotFroundUser
//	case err != nil:
//		return err
//	default:
//		log.Printf("update user with %v\n", filter.OfferID)
//	}
//
//	return nil
//}
