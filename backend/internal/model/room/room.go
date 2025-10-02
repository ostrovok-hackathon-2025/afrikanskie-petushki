package room

import "github.com/google/uuid"

type Room struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
