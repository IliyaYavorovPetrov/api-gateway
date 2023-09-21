package gateways

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Add(ctx context.Context, key string, val interface{}) error
	GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
