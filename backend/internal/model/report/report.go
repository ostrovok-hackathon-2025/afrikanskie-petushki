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
	ID           uuid.UUID
	UserID       uuid.UUID
	ExpirationAt time.Time
	Status       string
	Text         string
	Images       []Image
}
