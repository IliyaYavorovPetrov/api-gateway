package cache

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
)

func Sync(ctx context.Context, c1 gateways.Cache, c2 gateways.Cache) error {
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
