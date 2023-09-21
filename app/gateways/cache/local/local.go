package local

import (
	"context"
	"strings"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	icache "github.com/orcaman/concurrent-map/v2"
)

type Gateway struct {
	cache *icache.ConcurrentMap[string, interface{}]
}

var _ gateways.Cache = (*Gateway)(nil)
var instance *Gateway

func createInstance() *Gateway {
	c := icache.New[interface{}]()

	return &Gateway{
		cache: &c,
	}
}

func init() {
	instance = createInstance()
}

func GetInstance() *Gateway {
	return instance
}

func (gw *Gateway) Get(ctx context.Context, key string) (interface{}, error) {
	val, ok := gw.cache.Get(key)
	if !ok {
		return nil, cache.ErrNotFoundKey
	}

	return val, nil
}

func (gw *Gateway) Add(ctx context.Context, key string, val interface{}) error {
	gw.cache.Set(key, val)
	return nil
}

func (gw *Gateway) GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	var results []string
	keys := gw.cache.Keys()
	for _, item := range keys {
		if strings.HasPrefix(item, prefix) {
			results = append(results, item)
		}
	}

	return results, nil
}

func (gw *Gateway) Delete(ctx context.Context, key string) error {
	gw.cache.Pop(key)
	return nil
}

func (gw *Gateway) Flush(ctx context.Context) error {
	gw.cache.Clear()
	return nil
}
