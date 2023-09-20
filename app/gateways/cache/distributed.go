package cache

import (
	"context"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/redis/go-redis/v9"
)

func NewDistributedCache() *GatewayDistributed {
	return &GatewayDistributed{
		cache: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

var _ gateways.Cache = (*GatewayDistributed)(nil)

type GatewayDistributed struct {
	cache *redis.Client
}

func (gw *GatewayDistributed) Get(ctx context.Context, key string) (interface{}, error) {
	var obj interface{}
	err := gw.cache.HGetAll(ctx, key).Scan(&obj)
	if err != nil {
		return nil, ErrNotFoundKey
	}

	return "", nil
}

func (gw *GatewayDistributed) Set(ctx context.Context, key string, val interface{}) {
}

func (gw *GatewayDistributed) GetAll(ctx context.Context, key string) []interface{} {
	var results []interface{}
	return results
}

func (gw *GatewayDistributed) Delete(ctx context.Context, key string) {
}

func (gw *GatewayDistributed) Flush(ctx context.Context) {
}
