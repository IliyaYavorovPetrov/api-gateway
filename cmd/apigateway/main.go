package main

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	mw "github.com/IliyaYavorovPetrov/api-gateway/app/server/middleware"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/routing"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	router := mux.NewRouter()

	routing.Init(ctx)
	auth.Init(ctx)

	apiRoutes := router.PathPrefix("/api/v0").Subrouter()
	adminRoutes := router.PathPrefix("/admin/v0").Subrouter()

	apiRoutes.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
	})
	apiRoutes.Use(mw.Routing)
	//apiRoutes.Use(mw.Auth)
	//apiRoutes.Use(mw.RateLimitting)
	apiRoutes.Use(mw.Logger)
	apiRoutes.Use(mw.Transform)

	adminRoutes.HandleFunc("/routing/configuration", routing.AddRoutingCfgHandler).Methods(http.MethodPost)
	adminRoutes.HandleFunc("/routing/configuration/all", routing.GetAllRoutingCfgHandler).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
