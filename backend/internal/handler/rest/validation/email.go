package validation

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid email")
)

func ValidateEmail(email string) error {
	// TODO: add email validation
	return nil
}
