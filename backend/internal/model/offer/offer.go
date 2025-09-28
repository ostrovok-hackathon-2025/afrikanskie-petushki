package offer

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

type Offer struct {
	ID           uuid.UUID `db:"offer_id"`
	Task         string    `db:"task"`
	RoomID       uuid.UUID `db:"room_id"`
	RoomName     string    `db:"room_name"`
	HotelID      uuid.UUID `db:"hotel_id"`
	HotelName    string    `db:"hotel_name"`
	CheckIn      time.Time `db:"check_in_at"`
	CheckOut     time.Time `db:"check_out_at"`
	ExpirationAT time.Time `db:"expiration_at"`
}

type Filter struct {
	LocationID   uuid.UUID
	PageSettings PageSettings
}

type Create struct {
	Task         string
	RoomID       uuid.UUID
	HotelID      uuid.UUID
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
}

type Edit struct {
	OfferID      uuid.UUID
	Task         pkg.Opt[string]
	RoomID       pkg.Opt[uuid.UUID]
	HotelID      pkg.Opt[uuid.UUID]
	CheckIn      pkg.Opt[time.Time]
	CheckOut     pkg.Opt[time.Time]
	ExpirationAT pkg.Opt[time.Time]
}

type PageSettings struct {
	Limit  int
	Offset int
}
