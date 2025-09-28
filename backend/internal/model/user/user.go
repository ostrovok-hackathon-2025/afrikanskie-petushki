package user

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID
	OstrovokLogin string
	Email         string
	IsAdmin       bool
}
