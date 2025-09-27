package validation

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	ErrPasswordTooShort    = errors.New("password is too short")
	ErrPasswordTooLong     = errors.New("password is too long")
	ErrPasswordNoDigits    = errors.New("password must contain digits")
	ErrPasswordNoCapitals  = errors.New("password must contain capitals")
	ErrPasswordNoSpecChars = errors.New("password must contain special characters")
)

const _PASSWORD_MIN_LEN = 8
const _PASSWORD_MAX_LEN = 64

func ValidatePassword(password string) error {
	length := utf8.RuneCountInString(password)

	if length < _PASSWORD_MIN_LEN {
		return ErrPasswordTooShort
	}

	if length > _PASSWORD_MAX_LEN {
		return ErrPasswordTooLong
	}

	reDigits := regexp.MustCompile(`[0-9]`)

	if !reDigits.MatchString(password) {
		return ErrPasswordNoDigits
	}

	reCapitals := regexp.MustCompile(`[A-Z]`)

	if !reCapitals.MatchString(password) {
		return ErrPasswordNoCapitals
	}

	reSpecChars := regexp.MustCompile(`[*\-+!@#$%^&(),.?":{}|<>~\\\/\[\];'=_]`)

	if !reSpecChars.MatchString(password) {
		return ErrPasswordNoSpecChars
	}

	return nil

}
