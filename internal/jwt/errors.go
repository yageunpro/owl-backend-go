package jwt

import "errors"

var ErrNoSecret = errors.New("no secret key provided")
var ErrInvalidToken = errors.New("invalid token")
