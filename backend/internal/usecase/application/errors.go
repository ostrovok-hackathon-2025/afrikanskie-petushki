package application

import "errors"

var (
	ErrNotOwner = errors.New("application does not belong to user")
)
