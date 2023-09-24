package middleware

import (
	"log"
	"net/http"
)

func Transform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)

		next.ServeHTTP(w, r)
	})
}