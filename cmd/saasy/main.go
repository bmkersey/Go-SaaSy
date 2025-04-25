package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/config"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	conn, err := sql.Open("postgres", cfg.DB)
	if err != nil {
		log.Fatal("Could not establish connection to DB", err)
	}

	queries := db.New(conn)
	store := db.NewStore(queries)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SaaSy is running 🚀 (via Chi)"))
	})

	r.Post("/register", auth.RegisterHandler(store))
	r.Post("/login", auth.LoginHandler(store, cfg.JwtSecret))

	log.Printf("Starting server on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
