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
	loc = local.CreateInstance("test-cache")
	dist = distributed.New("test-cache")
}

func GetCtx() context.Context {
	return ctx
}
