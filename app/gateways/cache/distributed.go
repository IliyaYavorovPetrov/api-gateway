package cache

import (
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/redis/go-redis/v9"
)

func NewDistributedCache() *GatewayDistributed {
	return &GatewayDistributed{
		cache: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

var _ gateways.Cache = (*GatewayDistributed)(nil)

type GatewayDistributed struct {
	cache *redis.Client
}

func (gw *GatewayDistributed) Get(key string) (interface{}, bool) {
	return "", true
}

func (gw *GatewayDistributed) Set(key string, val interface{}) {
}

func (gw *GatewayDistributed) GetAll(key string) []interface{} {
	var results []interface{}
	return results
}

func (gw *GatewayDistributed) Delete(key string) {
}

func (gw *GatewayDistributed) Flush() {
}
