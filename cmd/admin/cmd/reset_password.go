package cmd

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var resetEmail string
var sendEmail bool

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Generate a password reset token for a user",
	Run: func(cmd *cobra.Command, args []string) {
		if resetEmail == "" {
			log.Fatal("Missing --email flag")
		}

		_ = godotenv.Load()
		dbURL := os.Getenv("DB_URL")
		conn, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal("Could not establish connection to DB", err)
		}
		defer conn.Close()
		queries := db.New(conn)
		store := db.NewStore(queries)

		user, err := store.GetUserByEmail(context.Background(), resetEmail)
		if err != nil {
			log.Fatalf("Error finding user: %v", err)
		}

		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			log.Fatalf("Error generating token: %v", err)
		}

		token := base64.URLEncoding.EncodeToString(b)
		exp := time.Now().Add(10 * time.Minute)

		err = store.CreatePasswordReset(context.Background(), db.CreatePasswordResetParams{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: exp,
		})
		if err != nil {
			log.Fatalf("Failed to insert reset token: %v", err)
		}

		fmt.Printf("✅ Password reset token created for %s\n", resetEmail)
		fmt.Printf("🔗 Reset link: https://localhost:8080/reset-password?token=%s\n", token)

		if sendEmail {
			sender := email.NewEmailSender()
			data := struct {
				Email string
				Link  string
			}{
				Email: user.Email,
				Link:  fmt.Sprintf("https://saasy.app/reset-password?token=%s", token),
			}

			err := sender.SendTemplate(user.Email, "Reset Your Password", "reset.html.tmpl", data)
			if err != nil {
				log.Fatalf("❌ Failed to send reset email: %v", err)
			}

			fmt.Println("📨 Reset email sent successfully")
		}
	},
}
