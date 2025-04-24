package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	db        string
	jwtSecret string
	port      string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env vars")
		os.Exit(1)
	}

	cfg := Config{
		db:        getEnv("DB_URL", ""),
		port:      getEnv("PORT", "8080"),
		jwtSecret: getEnv("JWT_SECRET", ""),
	}

	if cfg.db == "" || cfg.jwtSecret == "" {
		log.Fatal("Missing env variables. Check .env")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
