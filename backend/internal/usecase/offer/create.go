package offer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
)

var (
	ErrAlreadyExists = errors.New("offer: already exists")
)

func (u *useCase) Create(ctx context.Context, create *model.Create) (uuid.UUID, error) {
	check := model.Check{
		LocationID: create.LocationID,
		RoomTypeID: create.RoomTypeID,
		HotelID:    create.HotelID,
		CheckIn:    create.CheckIn,
		CheckOut:   create.CheckOut,
	}
	exists, err := u.repo.Check(ctx, &check)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, ErrAlreadyExists
	}

	return u.repo.Create(ctx, create)
}
