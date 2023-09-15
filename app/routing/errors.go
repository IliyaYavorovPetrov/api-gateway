package routing

import "errors"

var (
	ErrNotValidRequestKey        = errors.New("not valid request key")
	ErrNotValidCfgRoutingHashKey = errors.New("not valid routing configuration hash key")
)
