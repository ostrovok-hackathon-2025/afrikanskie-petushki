package pkg

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

var (
	alpha = 0.0149
	gamma = 0.17628
)

func ChooseByRating(users []model.User) uuid.UUID {
	contributions := make([]float64, len(users))

	sum := 0.0
	for i, user := range users {
		contribution := transformRatingToContribution(user.Rating, alpha, gamma)
		contributions[i] = contribution
		sum += contribution
	}

	target := rand.Float64() * sum
	for i, c := range contributions {
		target -= c
		if target < 0 {
			return users[i].ID
		}
	}
	return uuid.Nil
}

func transformRatingToContribution(rating int, alpha, gamma float64) float64 {
	// f(r): в диапазон [10,100)
	f := 100 - 90*math.Exp(-alpha*float64(rating))
	// g(r): степень
	return math.Pow(f, gamma)
}
