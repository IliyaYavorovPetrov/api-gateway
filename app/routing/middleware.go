package routing

import (
	"context"
	"log"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	// TODO: move all middlewares to their own package
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "IsAuthNeeded", false)

		log.Println("routing")
		rri, err := GetRoutingCfgFromMethodHTTPSourceURL(ctx, r.Method, r.Host)
		if err != nil {
			log.Fatalf("no such request in the routing configuration")
		}

		r.Host = rri.DestinationURL
		if (rri.IsAuthNeeded) {
			ctx = context.WithValue(r.Context(), "IsAuthNeeded", true)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
