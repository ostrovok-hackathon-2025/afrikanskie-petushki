package offer

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

type Offer struct {
	ID                uuid.UUID `db:"offer_id"`
	Task              string    `db:"task"`
	RoomID            uuid.UUID `db:"room_id"`
	RoomName          string    `db:"room_name"`
	HotelID           uuid.UUID `db:"hotel_id"`
	HotelName         string    `db:"hotel_name"`
	LocationID        uuid.UUID `db:"location_id"`
	LocationName      string    `db:"location_name"`
	CheckIn           time.Time `db:"check_in_at"`
	CheckOut          time.Time `db:"check_out_at"`
	ExpirationAt      time.Time `db:"expiration_at"`
	Status            string    `db:"status"`
	ParticipantsLimit uint      `db:"participants_limit"`
	ParticipantsCount uint      `db:"participants_count"`
}

type Filter struct {
	ID         pkg.Opt[uuid.UUID]
	LocationID pkg.Opt[uuid.UUID]
	Limit      uint64
	Offset     uint64
}

type Create struct {
	Task              string
	RoomID            uuid.UUID
	CheckIn           time.Time
	CheckOut          time.Time
	ExpirationAT      time.Time
	HotelID           uuid.UUID
	LocalID           uuid.UUID
	ParticipantsLimit uint
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
	Limit  uint64
	Offset uint64
}
