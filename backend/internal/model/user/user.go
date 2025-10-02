package user

import (
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/achievement"
)

type User struct {
	ID            uuid.UUID
	OstrovokLogin string
	Email         string
	IsAdmin       bool
	Rating        int // Сделал проверку на <0 в usecase
	Achievements  []achievement.Achievement
}
