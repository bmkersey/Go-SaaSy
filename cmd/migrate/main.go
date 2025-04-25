package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Missing DB_URL")
	}

	db, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		log.Fatal("goose: failed to open DB:", err)
	}
	defer db.Close()

	migrationsDir := "sql/schema"

	log.Println("Checking path:", migrationsDir)
	_, err = os.Stat(migrationsDir)
	if err != nil {
		log.Fatalf("Directory check failed: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatal("goose up failed:", err)
	}
}
