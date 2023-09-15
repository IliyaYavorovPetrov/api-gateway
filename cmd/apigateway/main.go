package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IliyaYavorovPetrov/api-gateway/app/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/logger"
	"github.com/IliyaYavorovPetrov/api-gateway/app/ratelimitting"
	"github.com/IliyaYavorovPetrov/api-gateway/app/routing"
	"github.com/IliyaYavorovPetrov/api-gateway/app/transform"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	apiRoutes := router.PathPrefix("/api/v0").Subrouter()
	adminRoutes := router.PathPrefix("/admin/v0").Subrouter()

	apiRoutes.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[ %s ] %s%s\n %v", r.Method, r.Host, r.URL.Path, r.Header)
	})
	apiRoutes.Use(routing.Middleware)
	apiRoutes.Use(auth.Middleware)
	apiRoutes.Use(ratelimitting.Middleware)
	apiRoutes.Use(logger.Middleware)
	apiRoutes.Use(transform.Middleware)

	adminRoutes.HandleFunc("/routing/configuration/all", routing.GetAllRoutingCfgHandler).Methods(http.MethodGet)
	adminRoutes.HandleFunc("/routing/configuration", routing.AddRoutingCfgHandler).Methods(http.MethodPost)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
