package local

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"strings"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	cacheprovider "github.com/orcaman/concurrent-map/v2"
)

type Gateway[V any] struct {
	name  string
	cache *cacheprovider.ConcurrentMap[string, V]
}

func New[V any](name string) gateways.Cache[V] {
	pool := cache.Pool()

	newCache := &Gateway[V]{
		name:  name,
		cache: &cacheprovider.ConcurrentMap[string, V]{},
	}

	(*pool)[name] = newCache

	return newCache
}

func (gw *Gateway[V]) Get(ctx context.Context, key string) (V, error) {
	var val V

	_, ok := gw.cache.Get(key)
	if !ok {
		return val, cache.ErrNotFoundKey
	}

	return val, nil
}

func (gw *Gateway[V]) Add(ctx context.Context, key string, val V) error {
	gw.cache.Set(key, val)
	return nil
}

func (gw *Gateway[V]) AddAllItems(ctx context.Context, other map[string]V) error {
	for key, val := range other {
		gw.cache.Set(key, val)
	}

	return nil
}

func (gw *Gateway[V]) GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	var results []string
	keys := gw.cache.Keys()
	for _, item := range keys {
		if strings.HasPrefix(item, prefix) {
			results = append(results, item)
		}
	}

	return results, nil
}

func (gw *Gateway[V]) GetAllItems(ctx context.Context) (map[string]V, error) {
	items := gw.cache.Items()
	return items, nil
}

func (gw *Gateway[V]) Delete(ctx context.Context, key string) error {
	gw.cache.Pop(key)
	return nil
}

func (gw *Gateway[V]) Flush(ctx context.Context) error {
	gw.cache.Clear()
	return nil
}
