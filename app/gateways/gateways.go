package gateways

import "context"

// TODO: Make the interface with generics https://stackoverflow.com/questions/71132124/how-to-solve-interface-method-must-have-no-type-parameters

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Add(ctx context.Context, key string, val interface{}) error
	AddAllItems(ctx context.Context, other map[string]interface{}) error
	GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error)
	GetAllItems(ctx context.Context) (map[string]interface{}, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
