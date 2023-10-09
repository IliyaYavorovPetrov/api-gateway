package auth

import "errors"

var (
	ErrNotValidSessionHashKey = errors.New("not valid session hash key")
)
