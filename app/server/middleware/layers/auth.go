package layers

import (
	"bytes"
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
			http.Error(w, "custom data not found", http.StatusInternalServerError)
			return
		}

		if !isAuthNeeded {
			var responseWriter ResponseCapture
			next.ServeHTTP(&responseWriter, r)

			if responseWriter.Status() == http.StatusFound {
				s := models.Session{
					UserID:        "some-id",
					Username:      "ivan",
					UserRole:      "User",
					IsBlacklisted: false,
				}

				id, err := auth.AddToSessionStore(r.Context(), s)
				if err != nil {
					return
				}

				log.Printf("new session was created with id %s\n", id)
			}

			// Copy the captured response to the original ResponseWriter
			responseWriter.CopyTo(w)
		} else {
			// Check for the session
			next.ServeHTTP(w, r)
		}
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
