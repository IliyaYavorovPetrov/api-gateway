package distributed

import (
	"context"
	"log"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	icache "github.com/redis/go-redis/v9"
)

type Gateway struct {
	cache *icache.Client
}

var _ gateways.Cache = (*Gateway)(nil)
var instance *Gateway

func init() {
	instance = createInstance()
}

func createInstance() *Gateway {
	return &Gateway{
		cache: icache.NewClient(&icache.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func GetInstance() *Gateway {
	return instance
}

func (gw *Gateway) Get(ctx context.Context, key string) (interface{}, error) {
	var val interface{}
	err := gw.cache.HGetAll(ctx, key).Scan(&val)
	if err != nil {
		return nil, cache.ErrNotFoundKey
	}

	return val, nil
}

func (gw *Gateway) Add(ctx context.Context, key string, val interface{}) error {
	if _, err := gw.cache.HSet(ctx, key, val).Result(); err != nil {
		return cache.ErrFailedToAdd
	}

	return nil
}

func (gw *Gateway) AddAllItems(ctx context.Context, other map[string]interface{}) error {
	_, err := gw.cache.Pipelined(ctx, func(pipe icache.Pipeliner) error {
		for key, val := range other {
			pipe.Set(ctx, key, val, -1)
		}

		return nil
	})

	if err != nil {
		log.Println("failed to execute the pipe")
		return err
	}

	return nil
}

func (gw *Gateway) GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	var results []string
	iter := gw.cache.Scan(ctx, 0, prefix, 0).Iterator()
	for iter.Next(ctx) {
		results = append(results, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (gw *Gateway) GetAllItems(ctx context.Context) (map[string]interface{}, error) {
	keys, err := gw.GetAllKeysByPrefix(ctx, "*")
	if err != nil {
		return nil, err
	}

	items := make(map[string]interface{})

	cmds, err := gw.cache.Pipelined(ctx, func(pipe icache.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(ctx, key)
		}

		return nil
	})
	if err != nil {
		log.Println("failed to execute the pipe")
		return nil, err
	}

	items = cmds[0].(*icache.MapStringInterfaceCmd).Val()

	return items, nil
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
