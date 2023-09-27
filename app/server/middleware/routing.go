package middleware

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/routing"
	"log"
	"net/http"
)

func Routing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), MiddlewareContextKey(IsAuthNeededKey), false)

		rri, err := routing.GetRoutingCfgFromMethodHTTPSourceURL(ctx, r.Method, r.Host)
		if err != nil {
			log.Fatalf("no such request in the routing configuration")
		}

		r.Host = rri.DestinationURL
		if rri.IsAuthNeeded {
			ctx = context.WithValue(r.Context(), MiddlewareContextKey(IsAuthNeededKey), true)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
