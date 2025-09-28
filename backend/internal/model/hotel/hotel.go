package hotel

import (
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/location"
)

type Hotel struct {
	ID       uuid.UUID `db:"hotel_id"`
	Name     string    `db:"hotel_name"`
	Location location.Location
}

type Create struct {
	Name       string
	LocationID uuid.UUID
}
