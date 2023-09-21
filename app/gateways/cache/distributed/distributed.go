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

	return "", nil
}

func (gw *Gateway) Add(ctx context.Context, key string, val interface{}) error {
	if _, err := gw.cache.HSet(ctx, key, val).Result(); err != nil {
		return cache.ErrFailedToAdd
	}

	return nil
}

func (gw *Gateway) GetAll(ctx context.Context, key string) []interface{} {
	var results []interface{}
	return results
}

func (gw *Gateway) Delete(ctx context.Context, key string) {
}

func (gw *Gateway) Flush(ctx context.Context) {
}
