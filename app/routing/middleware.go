package routing

import (
	"context"
	"log"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		log.Println("routing")
		rri, err := GetRoutingCfgFromMethodHTTPSourceURL(ctx, r.Method, r.Host)
		if err != nil {
			log.Fatalf("no such request in the routing configuration")
		}

		r.Host = rri.DestinationURL
		log.Println(r.Host)

		next.ServeHTTP(w, r)
	})
}
