package distributed

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/routing"
	icache "github.com/redis/go-redis/v9"
)

type Gateway struct {
	cache *icache.Client
}

var _ gateways.Cache = (*Gateway)(nil)
var instance *Gateway
var apiGatewayHash string

func init() {
	instance = createInstance("apiGatewayHash")
}

func createInstance(hash string) *Gateway {
	apiGatewayHash = hash
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
	data, err := gw.cache.HGet(ctx, apiGatewayHash, key).Result()
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

	if _, err := gw.cache.HSet(ctx, apiGatewayHash, key, string(serializedVal)).Result(); err != nil {
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
		fieldNames, nextCursor, err := gw.cache.HScan(ctx, apiGatewayHash, cursor, prefix+"*", 10).Result()
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
		if err != nil && err != cache.ErrNotFoundKey {
			return nil, err
		}

		if data != nil {
			items[key] = data
		}
	}

	return items, nil
}

func (gw *Gateway) Delete(ctx context.Context, key string) error {
	_, err := gw.cache.HDel(ctx, apiGatewayHash, key).Result()
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
