package cache

import "errors"

var (
	ErrNotFoundKey = errors.New("not found key")
	ErrFailedToAdd = errors.New("could not add an element")
	ErrUndefinedValueType = errors.New("unknown value type")
)
