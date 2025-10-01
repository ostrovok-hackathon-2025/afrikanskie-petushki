package report

import (
	"time"

	"github.com/google/uuid"
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
