package auth

import "errors"

var ErrInvalidToken = errors.New("Invalid token")
var ErrExpiredToken = errors.New("Token expired")
