package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/config"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	cfg := config.LoadConfig()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := sql.Open("postgres", cfg.DB)
	if err != nil {
		log.Fatal("Could not establish connection to DB", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	store := db.NewStore(queries)
	sender := email.NewEmailSender()

	emailAddr := "bmkersey@gmail.com"

	user, err := store.GetUserByEmail(context.Background(), emailAddr)
	if err != nil {
		log.Fatalf("User not found: %v", err)
	}

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		log.Fatalf("Random token generation failed: %v", err)
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(10 * time.Minute)

	err = store.CreatePasswordReset(context.Background(), db.CreatePasswordResetParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Fatalf("Failed to create pw reset: %v", err)
	}

	resetLink := fmt.Sprintf("https://saasy.app/reset-password?token=%s", token)

	data := struct {
		Email string
		Link  string
	}{
		Email: user.Email,
		Link:  resetLink,
	}

	err = sender.SendTemplate(emailAddr, "Reset Your SaaSy Password", "reset.html.tmpl", data)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("✅ Reset password email sent.")
}
