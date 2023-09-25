package cache

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
)

// If distributed write fails, we have inconsistent state of cache

func WriteToLocalAndDistributedCaches(ctx context.Context, local gateways.Cache, distributed gateways.Cache, key string, data interface{}) error {
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

// Again can have inconsistent state of cache c1

func LoadCacheOneWithCacheTwo(ctx context.Context, c1 gateways.Cache, c2 gateways.Cache) error {
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
