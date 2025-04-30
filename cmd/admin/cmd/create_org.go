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
	"github.com/spf13/cobra"
)

var (
	orgName string
	ownerID string
)

var createOrgCmd = &cobra.Command{
	Use:   "create-org",
	Short: "Create a new org with an owner",
	Run: func(cmd *cobra.Command, args []string) {
		if orgName == "" || ownerID == "" {
			log.Fatal("Both --name and --owner-id are required")
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

		ownerUUID, err := uuid.Parse(ownerID)
		if err != nil {
			log.Fatalf("invalid owner ID: %v", err)
		}

		orgID := uuid.New()

		_, err = store.CreateOrganization(context.Background(), db.CreateOrganizationParams{
			ID:      orgID,
			Name:    orgName,
			OwnerID: ownerUUID,
		})
		if err != nil {
			log.Fatalf("Failed to create organization: %v", err)
		}

		fmt.Printf("✅ Created organization: %s (ID: %s, Owner: %s)\n", orgName, orgID, ownerID)
	},
}
