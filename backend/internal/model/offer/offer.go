package offer

import (
	"time"

	"github.com/google/uuid"
)

type Offer struct {
	ID           uuid.UUID
	HotelID      uuid.UUID
	ExpirationAT time.Time
	LocalID      uuid.UUID
	Used         bool
	Task         string
}

type Filter struct {
	LocalID      uuid.UUID
	PageSettings PageSettings
}

type Create struct {
	HotelID      uuid.UUID
	ExpirationAT time.Time
	LocalID      uuid.UUID
	Task         string
}

type Edit struct {
	OfferID      uuid.UUID
	ExpirationAT time.Time
	Task         string
}

type PageSettings struct {
	Limit  int
	Offset int
}
