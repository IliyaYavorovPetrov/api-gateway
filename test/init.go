package test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
)

var ctx context.Context
var loc gateways.Cache
var dist gateways.Cache

func init() {
	ctx = context.Background()
	loc = local.GetInstance()
	dist = distributed.GetInstance()
}

func GetCtx() context.Context {
	return ctx
}

func GetLocalCache() gateways.Cache {
	return loc
}

func GetDistributedCache() gateways.Cache {
	return dist
}
