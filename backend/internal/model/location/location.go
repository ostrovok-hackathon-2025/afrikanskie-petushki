package location

import "github.com/google/uuid"

type Location struct {
	ID   uuid.UUID `db:"location_id"`
	Name string    `db:"location_name"`
}
