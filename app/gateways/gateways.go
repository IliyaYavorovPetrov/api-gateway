package gateways

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Add(ctx context.Context, key string, val interface{}) error
	GetAll(ctx context.Context, key string) []interface{}
	Delete(ctx context.Context, key string)
	Flush(ctx context.Context)
}
