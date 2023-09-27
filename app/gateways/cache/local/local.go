package local

import (
	"context"
	"strings"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	cacheprovider "github.com/orcaman/concurrent-map/v2"
)

var pool = make(map[string]Gateway)

type Gateway struct {
	name  string
	cache *cacheprovider.ConcurrentMap[string, interface{}]
}

var _ gateways.Cache = (*Gateway)(nil)

func CreateInstance(name string) *Gateway {
	c := cacheprovider.New[interface{}]()

	inst := Gateway{
		name:  name,
		cache: &c,
	}

	if _, ok := pool[name]; !ok {
		pool[name] = inst
	}

	return &inst
}

func GetInstance(name string) *Gateway {
	res, ok := pool[name]
	if ok != true {
		return nil
	}

	return &res
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

func (gw *Gateway) AddAllItems(ctx context.Context, other map[string]interface{}) error {
	for key, val := range other {
		gw.cache.Set(key, val)
	}

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

func (gw *Gateway) GetAllItems(ctx context.Context) (map[string]interface{}, error) {
	items := gw.cache.Items()
	return items, nil
}

func (gw *Gateway) Delete(ctx context.Context, key string) error {
	gw.cache.Pop(key)
	return nil
}

func (gw *Gateway) Flush(ctx context.Context) error {
	gw.cache.Clear()
	return nil
}
