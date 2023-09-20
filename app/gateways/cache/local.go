package cache

import (
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

func (gw *GatewayLocal) Get(key string) (interface{}, bool) {
	val, isFound := gw.cache.Get(key)
	if !isFound {
		return "", false
	}

	strVal, isStr := val.(string)
	if !isStr {
		return "", false
	}

	return strVal, true
}

func (gw *GatewayLocal) Set(key string, val interface{}) {
	gw.cache.Set(key, val, cache.DefaultExpiration)
}

func (gw *GatewayLocal) GetAll(key string) []interface{} {
	var results []interface{}

	for _, item := range gw.cache.Items() {
		results = append(results, item.Object)
	}

	return results
}

func (gw *GatewayLocal) Delete(key string) {
	gw.cache.Delete(key)
}

func (gw *GatewayLocal) Flush() {
	gw.cache.Flush()
}
