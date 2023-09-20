package cache

import (
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/patrickmn/go-cache"
)

func NewLocalCache() *Gateway {
	return &Gateway{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

var _ gateways.Cache = (*Gateway)(nil)

type Gateway struct {
	cache *cache.Cache
}

func (gw *Gateway) Get(key string) (interface{}, bool) {
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

func (gw *Gateway) Set(key string, val interface{}) {
	gw.cache.Set(key, val, cache.DefaultExpiration)
}

func (gw *Gateway) GetAll(key string) []interface{} {
	var results []interface{}

	for _, item := range gw.cache.Items() {
		results = append(results, item.Object)
	}

	return results
}

func (gw *Gateway) Delete(key string) {
	gw.cache.Delete(key)
}

func (gw *Gateway) Flush() {
	gw.cache.Flush()
}
