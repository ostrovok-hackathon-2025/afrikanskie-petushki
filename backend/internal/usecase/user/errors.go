package register

import "errors"

var ErrIncorrectPassword = errors.New("incorrect password")
var ErrInvalidToken = errors.New("invalid token")
var ErrUserExists = errors.New("user exists")
