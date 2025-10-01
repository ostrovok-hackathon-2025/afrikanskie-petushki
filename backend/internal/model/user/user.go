package user

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID
	OstrovokLogin string
	Email         string
	IsAdmin       bool
	Rating        int // Сделал проверку на <0 в usecase
}
