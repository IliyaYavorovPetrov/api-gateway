package auth

import "errors"

var (
	ErrNotValidUserRole       = errors.New("not valid user role")
	ErrNotValidSessionHashKey = errors.New("not valid session hash key")
)
