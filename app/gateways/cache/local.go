package cache

import (
	"context"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/patrickmn/go-cache"
)

func NewLocalCache() *GatewayLocal {
	return &GatewayLocal{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

var _ gateways.Cache = (*GatewayLocal)(nil)

type GatewayLocal struct {
	cache *cache.Cache
}

func (gw *GatewayLocal) Get(ctx context.Context, key string) (interface{}, error) {
	val, isFound := gw.cache.Get(key)
	if !isFound {
		return "", ErrNotFoundKey
	}

	strVal, isStr := val.(string)
	if !isStr {
		return "", ErrNotFoundKey
	}

	return strVal, nil
}

func (gw *GatewayLocal) Add(ctx context.Context, key string, val interface{}) error {
	gw.cache.Set(key, val, cache.DefaultExpiration)
	return nil
}

func (gw *GatewayLocal) GetAll(ctx context.Context, key string) []interface{} {
	var results []interface{}

	for _, item := range gw.cache.Items() {
		results = append(results, item.Object)
	}

	return results
}

func (gw *GatewayLocal) Delete(ctx context.Context, key string) {
	gw.cache.Delete(key)
}

func (gw *GatewayLocal) Flush(ctx context.Context) {
	gw.cache.Flush()
}
