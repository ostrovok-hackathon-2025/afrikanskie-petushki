package offer

import (
	"time"

	"github.com/google/uuid"
)

type Offer struct {
	ID           uuid.UUID
	Task         string
	Location     string
	RoomType     string
	Hotel        string
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
	Location     string
	RoomType     string
	Hotel        string
	CheckIn      time.Time
	CheckOut     time.Time
	ExpirationAT time.Time
	HotelID      uuid.UUID
	LocalID      uuid.UUID
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
	Location string
	RoomType string
	Hotel    string
	CheckIn  time.Time
	CheckOut time.Time
}
