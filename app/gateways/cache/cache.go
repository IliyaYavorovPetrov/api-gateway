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

func SyncFromTo[V any](ctx context.Context, from gateways.Cache[V], to gateways.Cache[V]) error {
	items, err := from.GetAllItems(ctx)
	if err != nil {
		return err
	}

	err = to.AddAllItems(ctx, items)
	if err != nil {
		return err
	}

	return nil
}
