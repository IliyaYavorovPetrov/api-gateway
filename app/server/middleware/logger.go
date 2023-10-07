package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func extractInfoRequest(r *http.Request) string {
	return fmt.Sprintf("%s %s %s %s", r.Method, r.Host, r.RequestURI, r.RemoteAddr)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(extractInfoRequest(r))

		next.ServeHTTP(w, r)
	})
}
