package application

import "errors"

var (
	ErrOfferNotExist       = errors.New("offer for this application not exists")
	ErrUserNotExist        = errors.New("user for this application not exists")
	ErrApplicationNotFound = errors.New("application not found")
	ErrPageNotFound        = errors.New("page not found")
	ErrParticipantsLimit   = errors.New("all places for participants taken")
	ErrAppLimit            = errors.New("user application limit reached")
)
