package gateways

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Add(ctx context.Context, key string, val interface{}) error
	AddAllItems(ctx context.Context, other map[string]interface{}) error
	GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error)
	GetAllItems(ctx context.Context) (map[string]interface{}, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
