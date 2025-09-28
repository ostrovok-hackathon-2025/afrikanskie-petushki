package hotel

import (
	"github.com/google/uuid"
)

type Hotel struct {
	ID           uuid.UUID `db:"hotel_id"`
	Name         string    `db:"hotel_name"`
	LocationID   uuid.UUID `db:"location_id"`
	LocationName string    `db:"location_name"`
}

type Create struct {
	Name       string
	LocationID uuid.UUID
}
