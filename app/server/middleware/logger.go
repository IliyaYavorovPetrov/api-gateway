package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func getUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func extractInfoRequest(r *http.Request) string {
	return fmt.Sprintf("%s %s %s %s", r.Method, r.Host, r.RequestURI, getUserIP(r))
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(extractInfoRequest(r))

		next.ServeHTTP(w, r)
	})
}
