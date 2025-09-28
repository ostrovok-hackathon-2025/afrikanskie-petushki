package offer

import (
	"time"

	"github.com/google/uuid"
)

type Hotel struct {
	ID   uuid.UUID
	Name string
}

type RoomType struct {
	ID   uuid.UUID
	Name string
}

type Location struct {
	ID   uuid.UUID
	Name string
}

type Offer struct {
	ID           uuid.UUID
	Task         string
	Location     Location
	RoomType     RoomType
	Hotel        Hotel
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
}

type Filter struct {
	LocalID      uuid.UUID
	PageSettings PageSettings
}

type Create struct {
	Task         string
	LocationID   uuid.UUID
	RoomTypeID   uuid.UUID
	HotelID      uuid.UUID
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
}

type Edit struct {
	OfferID      uuid.UUID
	Task         string
	Location     string
	RoomType     string
	Hotel        string
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
}

type PageSettings struct {
	Limit  int
	Offset int
}

type Check struct {
	LocationID uuid.UUID
	RoomTypeID uuid.UUID
	HotelID    uuid.UUID
	CheckIn    time.Time
	CheckOut   time.Time
}
