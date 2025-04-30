package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var (
	userEmail    string
	userPassword string
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user",
	Run: func(cmd *cobra.Command, args []string) {
		if userEmail == "" || userPassword == "" {
			log.Fatal("All fields email and password are required")
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

		hashedPW, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		id := uuid.New()
		_, err = store.CreateUser(context.Background(), db.CreateUserParams{
			Email:        userEmail,
			PasswordHash: string(hashedPW),
			ID:           id,
		})
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		fmt.Printf("✅ Created user: %s (%s)\n", userEmail, id)

	},
}
