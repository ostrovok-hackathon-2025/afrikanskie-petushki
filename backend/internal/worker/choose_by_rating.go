package worker

import (
	"math"
	"math/rand"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

var (
	alpha = 0.0149
	gamma = 0.17628
)

func chooseByRating(apps []*application.ApplicationWithRating) *application.ApplicationWithRating {
	contributions := make([]float64, len(apps))

	sum := 0.0
	for i, app := range apps {
		contribution := transformRatingToContribution(app.UserRating, alpha, gamma)
		contributions[i] = contribution
		sum += contribution
	}

	target := rand.Float64() * sum
	for i, c := range contributions {
		target -= c
		if target < 0 {
			return apps[i]
		}
	}
	return nil
}

func transformRatingToContribution(rating int, alpha, gamma float64) float64 {
	// f(r): в диапазон [10,100)
	f := 100 - 90*math.Exp(-alpha*float64(rating))
	// g(r): степень
	return math.Pow(f, gamma)
}
