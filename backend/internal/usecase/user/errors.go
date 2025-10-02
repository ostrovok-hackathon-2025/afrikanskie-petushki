package user

import "errors"

var ErrIncorrectPassword = errors.New("incorrect password")
var ErrInvalidToken = errors.New("invalid token")
