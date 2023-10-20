package layers

import (
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/middleware"
	"io"
	"log"
	"net/http"
)

func Transform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyReq, ok := r.Context().Value(middleware.ContextKey(middleware.ProxyRequest)).(*http.Request)
		if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create an HTTP client and send the request to the destination server.
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal("failed to close the request body")
			}
		}(resp.Body)

		// Copy the response status code and headers to the client's response.
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy the response body to the client's response.
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		next.ServeHTTP(w, r)
	})
}
