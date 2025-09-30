package pkg

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

var (
	pow = 0.01
)

func ChooseByRating(users []model.User) uuid.UUID {
	probabilities := make([]float64, len(users))

	sum := 0.0
	for i, user := range users {
		probability := probabilityFromRating(user.Rating)
		probabilities[i] = probability
		sum += probability
	}

	target := rand.Float64() * sum
	for i, p := range probabilities {
		target -= p
		if target < 0 {
			return users[i].ID
		}
	}
	return uuid.Nil
}

func probabilityFromRating(rating int) float64 {
	return math.Pow(float64(rating), pow)
}
