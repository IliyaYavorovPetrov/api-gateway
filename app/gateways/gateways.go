package gateways

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Add(ctx context.Context, key string, val interface{}) error
	GetAllKeysByPattern(ctx context.Context, pattern string) ([]interface{}, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
