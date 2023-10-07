package layers

import (
	"fmt"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/middleware"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthNeeded, ok := r.Context().Value(middleware.ContextKey(middleware.IsAuthNeededKey)).(bool)
		if !ok {
			http.Error(w, "custom data not found", http.StatusInternalServerError)
			return
		}

		if isAuthNeeded {
			fmt.Println("do sth")
		}

		next.ServeHTTP(w, r)
	})
}
