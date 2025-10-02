package worker

import (
	"fmt"
	"testing"
)

func TestContributionFromRating(t *testing.T) {

	var (
		alphas = []float64{0.0149}
		gammas = []float64{0.17628}
	)

	bigUserRating := 10_000
	normalUserRating := 100
	smallUserRating := 0
	smallUser2Rating := 10

	for _, alpha = range alphas {
		for _, gamma = range gammas {
			contributionOfBigUser := transformRatingToContribution(bigUserRating, alpha, gamma)
			contributionOfNormalUserRating := transformRatingToContribution(normalUserRating, alpha, gamma)
			contributionOfSmallUserRating := transformRatingToContribution(smallUserRating, alpha, gamma)
			contributionOfSmallUser2Rating := transformRatingToContribution(smallUser2Rating, alpha, gamma)
			fmt.Printf("With alpha = %+v, gamma= %+v\n", alpha, gamma)
			fmt.Println("Contribution of big user:    ", contributionOfBigUser)
			fmt.Println("Contribution of normal user: ", contributionOfNormalUserRating)
			fmt.Println("Contribution of small user:  ", contributionOfSmallUserRating)
			fmt.Println("Contribution of small user2: ", contributionOfSmallUser2Rating)
			fmt.Println("small if small vs normal user: ",
				contributionOfSmallUserRating/(contributionOfNormalUserRating+contributionOfSmallUserRating))
			fmt.Println("normal if normal vs big user: ",
				contributionOfNormalUserRating/(contributionOfNormalUserRating+contributionOfBigUser))
			fmt.Println("small if small vs small2 user: ",
				contributionOfSmallUserRating/(contributionOfSmallUserRating+contributionOfSmallUser2Rating))
			fmt.Println()
		}
	}

	//Result
	//With alpha = 0.0149, gamma= 0.17628
	//Contribution of big user:     2.251956514544276
	//Contribution of normal user:  2.1637390507268726
	//Contribution of small user:   1.5006520298004717
	//Contribution of small user2:  1.7995965185560954
	//small if small vs normal user:  0.40952289120420854
	//normal if normal vs big user:  0.4900109209847683
	//small if small vs small2 user:  0.4547087917204767

}
