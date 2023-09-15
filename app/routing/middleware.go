package routing

import (
	"context"
	"log"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		rri, err := GetRoutingCfgFromMethodHTTPSourceURL(ctx, r.Method, r.Host)
		if err != nil {
			log.Fatalf("no such request in the routing configuration")
		}

		r.Host = rri.DestinationURL

		if rri.IsAuthNeeded {
			next.ServeHTTP(w, r)
		}
	})
}
