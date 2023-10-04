package test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
	"log"
)

var ctx context.Context
var loc gateways.Cache[]
var dist gateways.Cache

func init() {
	ctx = context.Background()
	loc = local.New("test-cache")
	dist = distributed.New("test-cache")
}

func GetCtx() context.Context {
	return ctx
}

func ClearLocalCache() {
	err := loc.Flush(ctx)
	if err != nil {
		log.Fatal("failed to clear local storage")
	}
}

func ClearDistributedCache() {
	err := dist.Flush(ctx)
	if err != nil {
		log.Fatal("failed to clear local storage")
	}
}

func ClearCaches() {
	ClearLocalCache()
	ClearDistributedCache()
}

func ContainsItem(item string, arr []string) bool {
	for _, i := range arr {
		if i == item {
			return true
		}
	}

	return false
}

func MapEqual(a map[string]interface{}, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, ok := b[key]
		if !ok || valA != valB {
			return false
		}
	}

	return true
}
