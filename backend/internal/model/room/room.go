package room

import "github.com/google/uuid"

type Room struct {
	ID   uuid.UUID `db:"room_id"`
	Name string    `db:"room_name"`
}
