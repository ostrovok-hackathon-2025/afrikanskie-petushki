package offer

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/hotel"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/room"
)

type Offer struct {
	ID           uuid.UUID `db:"offer_id"`
	Task         string    `db:"task"`
	Room         room.Room
	Hotel        hotel.Hotel
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
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
	HotelID      uuid.UUID
	LocalID      uuid.UUID
}

type Edit struct {
	OfferID      uuid.UUID
	Task         string
	RoomID       uuid.UUID
	HotelID      uuid.UUID
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
}

type PageSettings struct {
	Limit  int
	Offset int
}

type Check struct {
	RoomTypeID uuid.UUID
	HotelID    uuid.UUID
	CheckIn    time.Time
	CheckOut   time.Time
}
