package distributed

import (
	"context"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	icache "github.com/redis/go-redis/v9"
)

type Gateway struct {
	cache *icache.Client
}

var _ gateways.Cache = (*Gateway)(nil)
var instance *Gateway

func createInstance() *Gateway {
	return &Gateway{
		cache: icache.NewClient(&icache.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func init() {
	instance = createInstance()
}

func GetInstance() *Gateway {
	return instance
}

func (gw *Gateway) Get(ctx context.Context, key string) (interface{}, error) {
	var obj interface{}
	err := gw.cache.HGetAll(ctx, key).Scan(&obj)
	if err != nil {
		return nil, cache.ErrNotFoundKey
	}

	return obj, nil
}

func (gw *Gateway) Add(ctx context.Context, key string, val interface{}) error {
	if _, err := gw.cache.HSet(ctx, key, val).Result(); err != nil {
		return cache.ErrFailedToAdd
	}

	return nil
}

func (gw *Gateway) GetAllKeysByPattern(ctx context.Context, pattern string) ([]interface{}, error) {
	var results []interface{}
	iter := gw.cache.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		results = append(results, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (gw *Gateway) Delete(ctx context.Context, key string) error {
	_, err := gw.cache.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (gw *Gateway) Flush(ctx context.Context) error {
	_, err := gw.cache.FlushAll(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
