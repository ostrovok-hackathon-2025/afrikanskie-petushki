package validation

func ValidateRating(rating int) int {
	if rating < 0 {
		rating = 0
	}
	return rating
}
