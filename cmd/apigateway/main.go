package main

import (
	"fmt"
	"github.com/IliyaYavorovPetrov/api-gateway/app/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/logger"
	"github.com/IliyaYavorovPetrov/api-gateway/app/ratelimitting"
	"github.com/IliyaYavorovPetrov/api-gateway/app/routing"
	"github.com/IliyaYavorovPetrov/api-gateway/app/transform"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[ %s ] %s%s\n %v", r.Method, r.Host, r.URL.Path, r.Header)
		w.WriteHeader(http.StatusOK)
	})
	r.Use(routing.Middleware)
	r.Use(auth.Middleware)
	r.Use(ratelimitting.Middleware)
	r.Use(logger.Middleware)
	r.Use(transform.Middleware)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
