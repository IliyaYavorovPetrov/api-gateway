package auth

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthNeeded, stat := r.Context().Value("IsAuthNeeded").(bool)
		if !stat {
            http.Error(w, "custom data not found", http.StatusInternalServerError)
            return
        }

		println(isAuthNeeded)

		next.ServeHTTP(w, r)
	})
}
