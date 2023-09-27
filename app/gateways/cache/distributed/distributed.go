package distributed

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/routing"
	cacheprovider "github.com/redis/go-redis/v9"
)

var cachePool = make(map[string]Gateway)

type Gateway struct {
	name  string
	cache *cacheprovider.Client
}

var _ gateways.Cache = (*Gateway)(nil)

func CreateInstance(name string) *Gateway {
	inst := Gateway{
		name: name,
		cache: cacheprovider.NewClient(&cacheprovider.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	if _, ok := cachePool[name]; !ok {
		cachePool[name] = inst
	}

	return &inst
}

func GetInstance(name string) *Gateway {
	res, ok := cachePool[name]
	if ok != true {
		return nil
	}

	return &res
}

func (gw *Gateway) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := gw.cache.HGet(ctx, gw.name, key).Result()
	if err != nil {
		return nil, cache.ErrNotFoundKey
	}

	var rri routing.ReqRoutingInfo
	if err := json.Unmarshal([]byte(data), &rri); err == nil {
		return rri, nil
	}

	var session auth.Session
	if err := json.Unmarshal([]byte(data), &session); err == nil {
		return session, nil
	}

	return nil, cache.ErrUndefinedValueType
}

func (gw *Gateway) Add(ctx context.Context, key string, val interface{}) error {
	serializedVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if _, err := gw.cache.HSet(ctx, gw.name, key, string(serializedVal)).Result(); err != nil {
		return err
	}

	return nil
}

func (gw *Gateway) AddAllItems(ctx context.Context, other map[string]interface{}) error {
	for key, val := range other {
		if err := gw.Add(ctx, key, val); err != nil {
			log.Printf("failed to add key '%s' to cache: %v", key, err)
			return err
		}
	}
	return nil
}

func (gw *Gateway) GetAllKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
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

func (gw *Gateway) GetAllItems(ctx context.Context) (map[string]interface{}, error) {
	keys, err := gw.GetAllKeysByPrefix(ctx, "")
	if err != nil {
		return nil, err
	}

	items := make(map[string]interface{})

	for _, key := range keys {
		data, err := gw.Get(ctx, key)
		if err != nil && !errors.Is(err, cache.ErrNotFoundKey) {
			return nil, err
		}

		if data != nil {
			items[key] = data
		}
	}

	return items, nil
}

func (gw *Gateway) Delete(ctx context.Context, key string) error {
	_, err := gw.cache.HDel(ctx, gw.name, key).Result()
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
