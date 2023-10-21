package layers

import (
	"bytes"
	"encoding/json"
	"github.com/IliyaYavorovPetrov/api-gateway/app/common/models"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/middleware"
	"io"
	"log"
	"net/http"
)

type ResponseCapture struct {
	status int
	body   bytes.Buffer
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthNeeded, ok := r.Context().Value(middleware.ContextKey(middleware.IsAuthNeededKey)).(bool)
		if !ok {
			http.Error(w, "auth data not found", http.StatusInternalServerError)
			return
		}

		proxyReq, ok := r.Context().Value(middleware.ContextKey(middleware.ProxyRequest)).(*http.Request)
		if !ok {
			http.Error(w, "auth data not found", http.StatusInternalServerError)
			return
		}

		// Did sth and now should receive from the service on :8081 session in the body of the req
		if !isAuthNeeded {
			var responseWriter ResponseCapture
			next.ServeHTTP(&responseWriter, r)

			var s models.Session
			// Decode the response body into the Session struct.
			err := json.NewDecoder(proxyReq.Body).Decode(&s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			id, err := auth.AddToSessionStore(r.Context(), s)
			if err != nil {
				return
			}

			log.Printf("new session was created with id %s\n", id)

			responseWriter.CopyTo(w)
			return
		}

		// Check in the session store
		next.ServeHTTP(w, r)
	})
}
func (rc *ResponseCapture) Header() http.Header {
	return http.Header{}
}

func (rc *ResponseCapture) Write(b []byte) (int, error) {
	rc.body.Write(b)
	return len(b), nil
}

func (rc *ResponseCapture) WriteHeader(statusCode int) {
	rc.status = statusCode
}

func (rc *ResponseCapture) Status() int {
	return rc.status
}

func (rc *ResponseCapture) CopyTo(w http.ResponseWriter) {
	w.WriteHeader(rc.status)
	_, err := io.Copy(w, &rc.body)
	if err != nil {
		return
	}
}
