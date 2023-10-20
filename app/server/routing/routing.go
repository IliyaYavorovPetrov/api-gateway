package routing

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/common/models"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
)

var delimiter = "|"
var prefixRoutingCfg = "routing:cfg:"

var localCache gateways.Cache[models.ReqRoutingInfo]
var distributedCache gateways.Cache[models.ReqRoutingInfo]

func Init(ctx context.Context) {
	crn := cron.New()
	localCache = local.New[models.ReqRoutingInfo]("routing-local-cache")
	distributedCache = distributed.New[models.ReqRoutingInfo]("routing-distributed-cache")

	_, err := crn.AddFunc("@every 10s", persistDistributedCache)
	if err != nil {
		panic("could not load start sync between local and distributed cache")
	}
	crn.Start()

	err = cache.SyncFromTo[models.ReqRoutingInfo](ctx, distributedCache, localCache)
	if err != nil {
		log.Fatal("could not load routing configuration")
	}
}

func persistDistributedCache() {
	ctx := context.Background()
	err := cache.SyncFromTo[models.ReqRoutingInfo](ctx, localCache, distributedCache)
	if err != nil {
		log.Fatal("could not load routing configuration")
	}
	log.Println("distributed cache is updated")
}

func CreateRoutingCfgHashKey(methodHTTP string, sourceURL string) string {
	var sb strings.Builder
	sb.WriteString(prefixRoutingCfg)
	sb.WriteString(methodHTTP)
	sb.WriteString(delimiter)
	sb.WriteString(sourceURL)
	return sb.String()
}

func ExtractRequestKeyFromRoutingCfgHashKey(s string) (string, error) {
	if strings.HasPrefix(s, prefixRoutingCfg) {
		return s[len(prefixRoutingCfg):], nil
	}

	return s, ErrNotValidCfgRoutingHashKey
}

func AddToRoutingCfgStore(ctx context.Context, rri models.ReqRoutingInfo) (string, error) {
	err := localCache.Add(ctx, CreateRoutingCfgHashKey(rri.MethodHTTP, rri.SourceURL), rri)
	if err != nil {
		log.Fatalf("failed to add a routing configuration %s", err)
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(rri.MethodHTTP)
	sb.WriteString(delimiter)
	sb.WriteString(rri.SourceURL)
	return sb.String(), nil
}

func GetRoutingCfgFromRequestKey(ctx context.Context, requestKey string) (models.ReqRoutingInfo, error) {
	rri, err := localCache.Get(ctx, requestKey)
	if err != nil {
		return models.ReqRoutingInfo{}, err
	}

	return *rri, nil
}

func GetAllRoutingCfgs(ctx context.Context) ([]string, error) {
	rris, err := localCache.GetAllKeysByPrefix(ctx, prefixRoutingCfg)
	if err != nil {
		return nil, err
	}

	return rris, nil
}

func RemoveRoutingCfgFromRoutingStore(ctx context.Context, requestKey string) error {
	err := localCache.Delete(ctx, requestKey)
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
		err = localCache.Delete(ctx, routingCfg)

		if err != nil {
			log.Fatalf("failed to delete a routing configuration %s", err)
		}
	}

	return nil
}
