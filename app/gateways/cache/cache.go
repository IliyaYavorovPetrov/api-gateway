package cache

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
)

var instances map[string]interface{}

func init() {
	instances = make(map[string]interface{})
}

func Pool() *map[string]interface{} {
	return &instances
}

func WriteToLocalAndDistributedCaches[V any](ctx context.Context, local gateways.Cache[V], distributed gateways.Cache[V], key string, data interface{}) error {
	err := local.Add(ctx, key, data)
	if err != nil {
		return err
	}

	err = distributed.Add(ctx, key, data)
	if err != nil {
		return err
	}

	return nil
}

func LoadCacheOneWithCacheTwo[V any](ctx context.Context, c1 gateways.Cache[V], c2 gateways.Cache[V]) error {
	items, err := c2.GetAllItems(ctx)
	if err != nil {
		return err
	}
	err = c1.AddAllItems(ctx, items)
	if err != nil {
		return err
	}

	return nil
}
