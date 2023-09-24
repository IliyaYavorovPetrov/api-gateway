package routing

import (
	"context"
	"encoding/json"
	"net/http"
)

func GetAllRoutingCfgHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cfgKeys, err := GetAllRoutingCfgs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var cfgs []ReqRoutingInfo
	
	for i := 0; i < len(cfgKeys); i++ {
        cfg, err := GetRoutingCfgFromRequestKey(ctx, cfgKeys[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cfgs = append(cfgs, cfg)
    }

	jcfgs, err := json.Marshal(cfgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    _, err = w.Write(jcfgs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func AddRoutingCfgHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var rri ReqRoutingInfo

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&rri); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	defer r.Body.Close()

	_, err := AddToRoutingCfgStore(ctx, &rri)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}