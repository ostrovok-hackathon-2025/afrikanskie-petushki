package report

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

type ImageURL string

type Image struct {
	ID   uuid.UUID
	Link string
}

type Report struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ApplicationID uuid.UUID
	ExpirationAt  time.Time
	Status        string
	Text          string
	HotelName     string
	LocationName  string
	Task          string
	CheckInAt     time.Time
	CheckOutAt    time.Time
	RoomName      string
	Images        []Image
}

func NewReport(applicationId uuid.UUID, ExpirationAt time.Time) Report {
	return Report{
		ID:            uuid.New(),
		ApplicationID: applicationId,
		ExpirationAt:  ExpirationAt,
		Status:        "created",
	}
}

type Filter struct {
	Status        pkg.Opt[string]
	HotelID       pkg.Opt[uuid.UUID]
	LocationID    pkg.Opt[uuid.UUID]
	Limit, Offset uint64
}
