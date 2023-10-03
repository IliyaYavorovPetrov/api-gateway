package gateways

import "context"

type Cache[V any] interface {
	Get(ctx context.Context, key string) (V, error)
	Add(ctx context.Context, key string, val V) error
	AddAllItems(ctx context.Context, other map[string]V) error
	GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error)
	GetAllItems(ctx context.Context) (map[string]V, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
