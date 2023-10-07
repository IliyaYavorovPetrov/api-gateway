package middleware

import (
	"net/http"
)

func Transform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.Host, http.StatusFound)

		next.ServeHTTP(w, r)
	})
}
