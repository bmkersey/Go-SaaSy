package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB                  string
	JwtSecret           string
	Port                string
	StripeWebhookSecret string
	AllowedPlans        map[string]string
	PlanPriceIDs        map[string]string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env vars")
		os.Exit(1)
	}

	cfg := Config{
		DB:                  getEnv("DB_URL", ""),
		Port:                getEnv("PORT", "8080"),
		JwtSecret:           getEnv("JWT_SECRET", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		AllowedPlans: map[string]string{
			getEnv("STRIPE_PRICE_ID_PRO", ""):      "pro",
			getEnv("STRIPE_PRICE_ID_ULTIMATE", ""): "ultimate",
		},
		PlanPriceIDs: map[string]string{
			"pro":      getEnv("STRIPE_PRICE_ID_PRO", ""),
			"ultimate": getEnv("STRIPE_PRICE_ID_ULTIMATE", ""),
		},
	}

	if cfg.DB == "" || cfg.JwtSecret == "" {
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
