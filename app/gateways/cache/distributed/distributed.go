package distributed

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"log"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	cacheprovider "github.com/redis/go-redis/v9"
)

type Gateway[V comparable] struct {
	name  string
	cache *cacheprovider.Client
}

func New[V comparable](name string) gateways.Cache[V] {
	pool := cache.Pool()

	newCache := &Gateway[V]{
		name: name,
		cache: cacheprovider.NewClient(&cacheprovider.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	(*pool)[name] = newCache

	return newCache
}

func (gw *Gateway[V]) Get(ctx context.Context, key string) (*V, error) {
	var val V

	data, err := gw.cache.HGet(ctx, gw.name, key).Result()
	if err != nil {
		return nil, cache.ErrNotFoundKey
	}

	if err := json.Unmarshal([]byte(data), &val); err != nil {
		return nil, cache.ErrNotFoundKey
	}

	return &val, nil
}

func (gw *Gateway[V]) Add(ctx context.Context, key string, val V) error {
	serializedVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if _, err := gw.cache.HSet(ctx, gw.name, key, string(serializedVal)).Result(); err != nil {
		return err
	}

	return nil
}

func (gw *Gateway[V]) AddAllItems(ctx context.Context, other map[string]V) error {
	for key, val := range other {
		if err := gw.Add(ctx, key, val); err != nil {
			log.Printf("failed to add key '%s' to cache: %v", key, err)
			return err
		}
	}
	return nil
}

func (gw *Gateway[V]) GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	var results []string

	cursor := uint64(0)
	isOdd := true

	for {
		fieldNames, nextCursor, err := gw.cache.HScan(ctx, gw.name, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, fieldName := range fieldNames {
			if isOdd {
				results = append(results, fieldName)
			}
			isOdd = !isOdd
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	return results, nil
}

func (gw *Gateway[V]) GetAllItems(ctx context.Context) (map[string]V, error) {
	keys, err := gw.GetAllKeysByPrefix(ctx, "")
	if err != nil {
		return nil, err
	}

	items := make(map[string]V)

	for _, key := range keys {
		data, err := gw.Get(ctx, key)
		if err != nil && !errors.Is(err, cache.ErrNotFoundKey) {
			return nil, err
		}

		if data != nil {
			items[key] = *data
		}
	}

	return items, nil
}

func (gw *Gateway[V]) Delete(ctx context.Context, key string) error {
	_, err := gw.cache.HDel(ctx, gw.name, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (gw *Gateway[V]) Flush(ctx context.Context) error {
	_, err := gw.cache.FlushAll(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
