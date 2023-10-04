package test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/models"
	"log"
)

var ctx context.Context
var loc gateways.Cache[models.ReqRoutingInfo]
var dist gateways.Cache[models.ReqRoutingInfo]

func init() {
	ctx = context.Background()
	loc = local.New[models.ReqRoutingInfo]("test-cache")
	dist = distributed.New[models.ReqRoutingInfo]("test-cache")
}

func GetCtx() context.Context {
	return ctx
}

func GetLocalCache() gateways.Cache[models.ReqRoutingInfo] {
	return loc
}

func GetWrongLocalCache() gateways.Cache[models.ReqRoutingInfo] {
	return &local.Gateway[models.ReqRoutingInfo]{}
}

func GetDistributedCache() gateways.Cache[models.ReqRoutingInfo] {
	return dist
}

func GetWrongDistributedCache() gateways.Cache[models.ReqRoutingInfo] {
	return &distributed.Gateway[models.ReqRoutingInfo]{}
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

func ClearBothCaches() {
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

func MapEqual[T comparable](a map[string]T, b map[string]T) bool {
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
