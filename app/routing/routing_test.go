package routing_test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/routing"
	"reflect"
	"testing"
)

func clearSessionStore(ctx context.Context) {
	routing.ClearRoutingCfgStore(ctx)
}

func TestAddAndGetFromRoutingCfgStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	rri1 := &routing.ReqRoutingInfo{
		SourceURL:      "https://src",
		DestinationURL: "http://dest",
		MethodHTTP:     "POST",
		IsAuthNeeded:   true,
	}

	reqKey, err := routing.AddToRoutingCfgStore(ctx, rri1)
	if err != nil {
		t.Fatalf("AddToRoutingCfgStore failed: %v", err)
	}

	methodHTTP, sourceURL, err := routing.GetMethodHTTPSourceURLFromRequestKey(reqKey)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rri2, err := routing.GetRoutingCfgFromMethodHTTPSourceURL(ctx, methodHTTP, sourceURL)
	if err != nil {
		t.Fatalf("GetRoutingCfgFromMethodHTTPSourceURL failed: %v", err)
	}

	if !reflect.DeepEqual(rri1, rri2) {
		t.Errorf("request session infos are different")
	}
}

func TestRemovingSessionFromSessionStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	rri1 := &routing.ReqRoutingInfo{
		SourceURL:      "https://src",
		DestinationURL: "http://dest",
		MethodHTTP:     "POST",
		IsAuthNeeded:   true,
	}

	reqKey1, err := routing.AddToRoutingCfgStore(ctx, rri1)
	if err != nil {
		t.Fatalf("AddToRoutingCfgStore failed: %v", err)
	}

	rri2 := &routing.ReqRoutingInfo{
		SourceURL:      "https://foo",
		DestinationURL: "http://bar",
		MethodHTTP:     "GET",
		IsAuthNeeded:   false,
	}

	reqKey2, err := routing.AddToRoutingCfgStore(ctx, rri2)
	if err != nil {
		t.Fatalf("AddToRoutingCfgStore failed: %v", err)
	}

	valuesToCheck := make(map[string]struct{})
	valuesToCheck[reqKey1] = struct{}{}
	valuesToCheck[reqKey2] = struct{}{}

	allReqRoutingInfos, err := routing.GetAllRoutingCfgs(ctx)
	if err != nil {
		t.Fatalf("GetAllRoutingCfgs failed: %v", err)
	}

	if len(allReqRoutingInfos) != len(valuesToCheck) {
		t.Errorf("wrong number of request infos, %d expected, %d received", len(allReqRoutingInfos), len(valuesToCheck))
	}

	for _, str := range allReqRoutingInfos {
		sID, err := routing.GetRequestKeyFromRoutingCfgHashKey(str)
		if err != nil {
			t.Errorf("%v", err)
		}
		if _, found := valuesToCheck[sID]; !found {
			t.Errorf("value %s not found in the list", str)
		}
	}
}
