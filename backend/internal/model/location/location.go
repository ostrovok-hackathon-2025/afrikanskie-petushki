package location

import "github.com/google/uuid"

type Location struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
