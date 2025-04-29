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
)

func WebhookHandler(store db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside webhook handler")
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Request body read error", http.StatusServiceUnavailable)
			return
		}
		fmt.Println(payload)

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		if endpointSecret == "" {
			http.Error(w, "Webhook secret missing", http.StatusInternalServerError)
			return
		}
		fmt.Println(r.Header.Get("Stripe-Signature"))
		fmt.Println(endpointSecret)
		// event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
		// if err != nil {
		// 	http.Error(w, "Webhook verification failed", http.StatusBadRequest)
		// 	return
		// }

		event := stripe.Event{}

		if err := json.Unmarshal(payload, &event); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(event)

		switch event.Type {
		case "checkout.session.completed":
			fmt.Println("Inside checkout sysytem completed")
			var session stripe.CheckoutSession
			if err := json.NewDecoder(bytes.NewReader(event.Data.Raw)).Decode(&session); err != nil {
				http.Error(w, "Webhook decode failed", http.StatusBadRequest)
				return
			}
			fmt.Printf("CheckoutSession: %+v\n", session)

			orgID := session.ClientReferenceID
			if orgID == "" {
				http.Error(w, "Missing org ID in session", http.StatusBadRequest)
				return
			}
			fmt.Println(orgID)
			// subscriptionID := session.Subscription
			// customerID := session.Customer

			orgUUID, err := uuid.Parse(orgID)
			if err != nil {
				http.Error(w, "Error parsing ID to UUID", http.StatusInternalServerError)
				return
			}

			priceID := ""
			if len(session.LineItems.Data) > 0 {
				priceID = session.LineItems.Data[0].Price.ID
			}

			plan, ok := allowedPlans[priceID]
			if !ok {
				plan = "basic" // fallback if unknown
			}
			fmt.Println(plan)
			stripeCustomerID := ""
			if session.Customer != nil {
				stripeCustomerID = session.Customer.ID
			}

			stripeSubscriptionID := ""
			if session.Subscription != nil {
				stripeSubscriptionID = session.Subscription.ID
			}

			if stripeCustomerID == "" || stripeSubscriptionID == "" {
				http.Error(w, "Missing customer or subscription ID in session", http.StatusBadRequest)
				return
			}

			// safe to continue
			err = store.UpdateOrganizationBilling(r.Context(), db.UpdateOrganizationBillingParams{
				ID: orgUUID,
				StripeCustomerID: sql.NullString{
					String: stripeCustomerID,
					Valid:  true,
				},
				StripeSubscriptionID: sql.NullString{
					String: stripeSubscriptionID,
					Valid:  true,
				},
				Plan: sql.NullString{
					String: plan,
					Valid:  true,
				},
				IsPaid: true,
			})
			if err != nil {
				fmt.Println("Error updating org billing", err)
				http.Error(w, "Failed to update billing", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

		case "customer.subscription.deleted":
			fmt.Println("Inside customer.subscription.deleted")
			var subscription stripe.Subscription
			if err := json.NewDecoder(bytes.NewReader(event.Data.Raw)).Decode(&subscription); err != nil {
				http.Error(w, "Webhook decode failer", http.StatusBadRequest)
				return
			}

			orgID := subscription.Metadata["org_id"]
			if orgID == "" {
				http.Error(w, "Missing org ID in subscription metadata", http.StatusBadRequest)
				return
			}

			orgUUID, err := uuid.Parse(orgID)
			if err != nil {
				http.Error(w, "Error parsing ID to UUID", http.StatusInternalServerError)
				return
			}

			priceID := ""
			if len(subscription.Items.Data) > 0 {
				priceID = subscription.Items.Data[0].Price.ID
			}

			plan, ok := allowedPlans[priceID]
			if !ok {
				plan = "free"
			}

			err = store.UpdateOrganizationBilling(r.Context(), db.UpdateOrganizationBillingParams{
				ID: orgUUID,
				StripeCustomerID: sql.NullString{
					String: "",
					Valid:  true,
				},
				StripeSubscriptionID: sql.NullString{
					String: "",
					Valid:  true,
				},
				Plan: sql.NullString{
					String: plan,
					Valid:  true,
				},
				IsPaid: false,
			})
			if err != nil {
				fmt.Println("Error downgrading org: ", err)
				http.Error(w, "Failed to downgrade org", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		case "customer.subscription.updated":
			fmt.Println("Inside customer sub updated")
			var subscription stripe.Subscription
			if err := json.NewDecoder(bytes.NewReader(event.Data.Raw)).Decode(&subscription); err != nil {
				http.Error(w, "Error decoding event data", http.StatusBadRequest)
				return
			}

			orgID := subscription.Metadata["org_id"]
			if orgID == "" {
				http.Error(w, "Missing org ID in subscription metadata", http.StatusBadRequest)
				return
			}

			orgUUID, err := uuid.Parse(orgID)
			if err != nil {
				http.Error(w, "Error parsing ID to UUID", http.StatusInternalServerError)
				return
			}

			priceID := ""
			if len(subscription.Items.Data) > 0 {
				priceID = subscription.Items.Data[0].Price.ID
			}

			plan, ok := allowedPlans[priceID]
			if !ok {
				plan = "free"
			}
			if subscription.Status == "active" {
				err = store.UpdateOrganizationBilling(r.Context(), db.UpdateOrganizationBillingParams{
					ID: orgUUID,
					StripeCustomerID: sql.NullString{
						String: subscription.Customer.ID,
						Valid:  true,
					},
					StripeSubscriptionID: sql.NullString{
						String: subscription.ID,
						Valid:  true,
					},
					Plan: sql.NullString{
						String: plan,
						Valid:  true,
					},
					IsPaid: false,
				})
				if err != nil {
					fmt.Println("Error downgrading organization: ", err)
					http.Error(w, "Failed to downgrade organization", http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusOK)
		}
	}
}
