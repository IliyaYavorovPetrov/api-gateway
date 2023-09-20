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

func (gw *Gateway) Get(key string) (string, bool) {
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

func (gw *Gateway) Set(key string, val string) {
	gw.cache.Set(key, val, cache.DefaultExpiration)
}

func (gw *Gateway) GetAll(key string) []string {
	var results []string

	for _, item := range gw.cache.Items() {
		results = append(results, item.Object.(string))
	}

	return results
}

func (gw *Gateway) Delete(key string) {
	gw.cache.Delete(key)
}

func (gw *Gateway) Flush() {
	gw.cache.Flush()
}
