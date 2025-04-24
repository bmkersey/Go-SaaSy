package main

import (
	"log"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SaaSy is running 🚀 (via Chi)"))
	})

	log.Printf("Starting server on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
