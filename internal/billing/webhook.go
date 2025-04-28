package billing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

func WebhookHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Request body read error", http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_SECRET")
		if endpointSecret == "" {
			http.Error(w, "Webhook secret missing", http.StatusInternalServerError)
			return
		}

		event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
		if err != nil {
			http.Error(w, "Webhook verification failed", http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			var session stripe.CheckoutSession
			if err := json.NewDecoder(bytes.NewReader(event.Data.Raw)).Decode(&session); err != nil {
				http.Error(w, "Webhook decode failed", http.StatusBadRequest)
				return
			}

			orgID := session.ClientReferenceID
			if orgID == "" {
				http.Error(w, "Missing org ID in session", http.StatusBadRequest)
				return
			}

			subscriptionID := session.Subscription.ID
			customerID := session.Customer.ID

			orgUUID, err := uuid.Parse(orgID)
			if err != nil {
				http.Error(w, "Error parsing ID to UUID", http.StatusInternalServerError)
				return
			}

			err = store.UpdateOrganizationBilling(r.Context(), db.UpdateOrganizationBillingParams{
				ID: orgUUID,
				StripeCustomerID: sql.NullString{
					String: customerID,
					Valid:  true,
				},
				StripeSubscriptionID: sql.NullString{
					String: subscriptionID,
					Valid:  true,
				},
				Plan: sql.NullString{
					String: "PLACE HOLDER",
					Valid:  true,
				},
			})
			if err != nil {
				fmt.Println("Error updating org billing", err)
				http.Error(w, "Failed to update billing", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusOK)
		}
	}
}
