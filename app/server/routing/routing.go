package routing

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/models"
	"log"
	"strings"
)

var delimiter = "|"
var prefixRoutingCfg = "routing:cfg:"

var localCache gateways.Cache[models.ReqRoutingInfo]
var distributedCache gateways.Cache[models.ReqRoutingInfo]

func init() {
	localCache = local.New[models.ReqRoutingInfo]("routing-local-cache")
	distributedCache = distributed.New[models.ReqRoutingInfo]("routing-distributed-cache")
}

func createRoutingCfgHashKey(methodHTTP string, sourceURL string) string {
	return prefixRoutingCfg + methodHTTP + delimiter + sourceURL
}

func ExtractRequestKeyFromRoutingCfgHashKey(s string) (string, error) {
	if strings.HasPrefix(s, prefixRoutingCfg) {
		return s[len(prefixRoutingCfg):], nil
	}

	return s, ErrNotValidCfgRoutingHashKey
}

func AddToRoutingCfgStore(ctx context.Context, rri models.ReqRoutingInfo) (string, error) {
	err := distributedCache.Add(ctx, createRoutingCfgHashKey(rri.MethodHTTP, rri.SourceURL), rri)
	if err != nil {
		log.Fatalf("failed to add a routing configuration %s", err)
		return "", err
	}

	return rri.MethodHTTP + delimiter + rri.SourceURL, nil
}

func GetRoutingCfgFromRequestKey(ctx context.Context, requestKey string) (models.ReqRoutingInfo, error) {
	rri, err := distributedCache.Get(ctx, prefixRoutingCfg+requestKey)
	if err != nil {
		return models.ReqRoutingInfo{}, err
	}

	return *rri, nil
}

func GetAllRoutingCfgs(ctx context.Context) ([]string, error) {
	rris, err := distributedCache.GetAllKeysByPrefix(ctx, prefixRoutingCfg)
	if err != nil {
		return nil, err
	}

	return rris, nil
}

func RemoveRoutingCfgFromRoutingStore(ctx context.Context, requestKey string) error {
	err := distributedCache.Delete(ctx, requestKey)
	if err != nil {
		return err
	}

	return nil
}

func ClearRoutingCfgStore(ctx context.Context) error {
	routingCfgs, err := GetAllRoutingCfgs(ctx)
	if err != nil {
		return err
	}

	for _, routingCfg := range routingCfgs {
		err = distributedCache.Delete(ctx, routingCfg)

		if err != nil {
			log.Fatalf("failed to delete a routing configuration %s", err)
		}
	}

	return nil
}
