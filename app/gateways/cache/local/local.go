package local

import (
	"context"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	icache "github.com/patrickmn/go-cache"
)

type Gateway struct {
	cache *icache.Cache
}

var _ gateways.Cache = (*Gateway)(nil)
var instance *Gateway

func createInstance() *Gateway {
	return &Gateway{
		cache: icache.New(icache.NoExpiration, icache.NoExpiration),
	}
}

func init() {
	instance = createInstance()
}

func GetInstance() *Gateway {
	return instance
}
	
func (gw *Gateway) Get(ctx context.Context, key string) (interface{}, error) {
	val, isFound := gw.cache.Get(key)
	if !isFound {
		return "", cache.ErrNotFoundKey
	}

	strVal, isStr := val.(string)
	if !isStr {
		return "", cache.ErrNotFoundKey
	}

	return strVal, nil
}

func (gw *Gateway) Add(ctx context.Context, key string, val interface{}) error {
	gw.cache.Set(key, val, icache.DefaultExpiration)
	return nil
}

func (gw *Gateway) GetAllKeysByPattern(ctx context.Context, pattern string) ([]interface{}, error) {
	var results []interface{}
	return results, nil
}

func (gw *Gateway) Delete(ctx context.Context, key string) error {
	gw.cache.Delete(key)
	return nil
}

func (gw *Gateway) Flush(ctx context.Context) error{
	gw.cache.Flush()
	return nil
}
