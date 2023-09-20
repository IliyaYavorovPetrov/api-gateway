package cache

import (
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/redis/go-redis/v9"
)

func NewDistributedCache() *GatewayLocal {
	return &GatewayLocal{
		cache: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

var _ gateways.Cache = (*GatewayLocal)(nil)

type GatewayLocal struct {
	cache *redis.Client
}

func (gw *GatewayLocal) Get(key string) (interface{}, bool) {
	return "", true
}

func (gw *GatewayLocal) Set(key string, val interface{}) {
}

func (gw *GatewayLocal) GetAll(key string) []interface{} {
	var results []interface{}
	return results
}

func (gw *GatewayLocal) Delete(key string) {
}

func (gw *GatewayLocal) Flush() {
}
