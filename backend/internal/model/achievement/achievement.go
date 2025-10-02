package achievement

import "github.com/google/uuid"

type Achievement struct {
	Id           uuid.UUID
	Name         string
	RaitingLimit int
}
