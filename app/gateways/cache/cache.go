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

func LoadInfo[V any](ctx context.Context, local gateways.Cache[V], distributed gateways.Cache[V]) error {
	items, err := distributed.GetAllItems(ctx)
	if err != nil {
		return err
	}

	err = local.AddAllItems(ctx, items)
	if err != nil {
		return err
	}

	return nil
}
