package main

import (
	"fmt"
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
