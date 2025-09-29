package hotel

import (
	"github.com/google/uuid"
)

type Hotel struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	LocationID   uuid.UUID `db:"location_id"`
	LocationName string    `db:"location_name"`
}

type Create struct {
	Name       string
	LocationID uuid.UUID
}
