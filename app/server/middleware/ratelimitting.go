package middleware

import (
	"log"
	"net/http"
)

func RateLimitting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host)

		next.ServeHTTP(w, r)
	})
}