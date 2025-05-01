package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var makeAdminEmail string

var makeAdminCmd = &cobra.Command{
	Use:   "make-admin",
	Short: "Promote a user to admin",
	Run: func(cmd *cobra.Command, args []string) {
		if makeAdminEmail == "" {
			log.Fatal("You must provide --email")
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

		err = store.SetUserAdmin(context.Background(), makeAdminEmail)
		if err != nil {
			log.Fatalf("Failed to promote user: %v", err)
		}

		fmt.Printf("✅ User %s is now an admin\n", makeAdminEmail)
	},
}

func init() {
	rootCmd.AddCommand(makeAdminCmd)
	makeAdminCmd.Flags().StringVar(&makeAdminEmail, "email", "", "User email to promote")
}
