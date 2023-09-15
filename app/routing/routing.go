package routing

import (
	"context"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var delimiter = "|"
var prefixRoutingCfg = "routing:cfg:"

func createRequestKey(methodHTTP string, sourceURL string) string {
	return methodHTTP + delimiter + sourceURL
}

func GetMethodHTTPSourceURLFromRequestKey(s string) (string, string, error) {
	ps := strings.Split(s, delimiter)

	if len(ps) == 2 {
		methodHTTP := strings.TrimSpace(ps[0])
		sourceURL := strings.TrimSpace(ps[1])
		return methodHTTP, sourceURL, nil
	}

	return "", "", ErrNotValidRequestKey
}

func createRoutingCfgHashKey(methodHTTP string, sourceURL string) string {
	return prefixRoutingCfg + createRequestKey(methodHTTP, sourceURL)
}

func GetRequestKeyFromRoutingCfgHashKey(s string) (string, error) {
	if strings.HasPrefix(s, prefixRoutingCfg) {
		return s[len(prefixRoutingCfg):], nil
	}

	return s, ErrNotValidCfgRoutingHashKey
}

func AddToRoutingCfgStore(ctx context.Context, rri *ReqRoutingInfo) (string, error) {
	if _, err := rdb.HSet(ctx, createRoutingCfgHashKey(rri.MethodHTTP, rri.SourceURL), rri).Result(); err != nil {
		log.Fatalf("failed to add a routing configuration %s", err)
		return "", err
	}

	return createRequestKey(rri.MethodHTTP, rri.SourceURL), nil
}

func GetRoutingCfgFromMethodHTTPSourceURL(ctx context.Context, methodHTTP string, sourceURL string) (ReqRoutingInfo, error) {
	rri, err := GetRoutingCfgFromRequestKey(ctx, createRoutingCfgHashKey(methodHTTP, sourceURL))
	if err != nil {
		return ReqRoutingInfo{}, err
	}

	return rri, nil
}

func GetRoutingCfgFromRequestKey(ctx context.Context, requestKey string) (ReqRoutingInfo, error) {
	var rri ReqRoutingInfo
	err := rdb.HGetAll(ctx, requestKey).Scan(&rri)
	if err != nil {
		return ReqRoutingInfo{}, err
	}

	return rri, nil
}

func GetAllRoutingCfgs(ctx context.Context) ([]string, error) {
	var rri []string
	iter := rdb.Scan(ctx, 0, prefixRoutingCfg+"*", 0).Iterator()
	for iter.Next(ctx) {
		rri = append(rri, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return rri, nil
}

func ClearRoutingCfgStore(ctx context.Context) {
	// TODO: Delete by pattern
	rdb.FlushAll(ctx)
}
