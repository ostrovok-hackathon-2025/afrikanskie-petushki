package validation

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	ErrUsernameTooShort = errors.New("username is too short")
	ErrUsernameTooLong  = errors.New("username is too long")
	ErrUsernameBadChar  = errors.New("username can only contain latin, digits and _")
)

const _USERNAME_MIN_LEN = 3
const _USERNAME_MAX_LEN = 32

func ValidateUsername(username string) error {
	length := utf8.RuneCountInString(username)

	if length < _USERNAME_MIN_LEN {
		return ErrUsernameTooShort
	}

	if length > _USERNAME_MAX_LEN {
		return ErrUsernameTooLong
	}

	reLatinDigits := regexp.MustCompile(`[\W]`)

	if reLatinDigits.MatchString(username) {
		return ErrUsernameBadChar
	}

	return nil
}
