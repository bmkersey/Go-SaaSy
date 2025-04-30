package billing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"
	"os"

	"github.com/bmkersey/Go-SaaSy/internal/config"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	stripeSession "github.com/stripe/stripe-go/v82/checkout/session"
	stripeSubscription "github.com/stripe/stripe-go/v82/subscription"
)

func WebhookHandler(store db.Store, sender *email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside webhook handler")
		cfg := config.LoadConfig()
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Request body read error", http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		if endpointSecret == "" {
			http.Error(w, "Webhook secret missing", http.StatusInternalServerError)
			return
		}

		event := stripe.Event{}

		if err := json.Unmarshal(payload, &event); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			fmt.Println("Inside checkout sysytem completed")
			var session stripe.CheckoutSession
			if err := json.NewDecoder(bytes.NewReader(event.Data.Raw)).Decode(&session); err != nil {
				http.Error(w, "Webhook decode failed", http.StatusBadRequest)
				return
			}
			fmt.Printf("CheckoutSession: %+v\n", session)

			params := &stripe.CheckoutSessionParams{
				Expand: []*string{stripe.String("line_items")},
			}
			fullSession, err := stripeSession.Get(session.ID, params)
			if err != nil {
				http.Error(w, "Failed to retrieve session: "+err.Error(), http.StatusInternalServerError)
				return
			}
			var priceID string
			// unpack the line items?
			if fullSession.LineItems != nil && len(fullSession.LineItems.Data) > 0 {
				for i, lineItem := range fullSession.LineItems.Data {
					if i == 0 {
						priceID = lineItem.Price.ID
						fmt.Println("Price ID:", priceID)
					}
				}
			}

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
			fmt.Println("---- DEBUG: using cfg.AllowedPlans map ----")
			for k, v := range cfg.AllowedPlans {
				fmt.Printf("Key: %s → Value: %s\n", k, v)
			}
			fmt.Printf("Looking up priceID: %s\n", priceID)

			plan, ok := cfg.AllowedPlans[priceID]
			fmt.Printf("Resolved plan: %s (found: %v)\n", plan, ok)

			if !ok {
				plan = "basic"
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
			org, err := store.GetOrganization(r.Context(), orgUUID)
			if err != nil {
				http.Error(w, "Error fetching Organization", http.StatusBadRequest)
				return
			}

			user, err := store.GetUserByID(r.Context(), org.OwnerID)
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				return
			}

			go func() {
				data := struct {
					Email string
					Plan  string
				}{
					Email: user.Email,
					Plan:  plan,
				}

				err := sender.SendTemplate(user.Email, "Your SaaSy Plan Changed", "plan_change.html.tmpl", data)
				if err != nil {
					log.Printf("Failed to send plan change email: %v", err)
				}
				fmt.Println("Email sent successfully")
			}()

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

			plan, ok := cfg.AllowedPlans[priceID]
			if !ok {
				plan = "basic"
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

			org, err := store.GetOrganization(r.Context(), orgUUID)
			if err != nil {
				http.Error(w, "Error fetching Organization", http.StatusBadRequest)
				return
			}

			user, err := store.GetUserByID(r.Context(), org.OwnerID)
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				return
			}

			go func() {
				data := struct {
					Email string
					Plan  string
				}{
					Email: user.Email,
					Plan:  plan,
				}

				err := sender.SendTemplate(user.Email, "Your SaaSy Plan Changed", "plan_change.html.tmpl", data)
				if err != nil {
					log.Printf("Failed to send plan change email: %v", err)
				}
			}()

			w.WriteHeader(http.StatusOK)
		case "customer.subscription.updated":
			fmt.Println("Inside customer sub updated")

			var evtSub stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &evtSub); err != nil {
				http.Error(w, "Webhook decode failed", http.StatusBadRequest)
				return
			}

			fullSub, err := stripeSubscription.Get(
				evtSub.ID,
				&stripe.SubscriptionParams{
					Params: stripe.Params{
						Expand: []*string{
							stripe.String("items.data.price"),
						},
					},
				},
			)
			if err != nil {
				http.Error(w, "Failed to fetch subscription", http.StatusInternalServerError)
				return
			}

			if len(fullSub.Items.Data) == 0 {
				http.Error(w, "No subscription items found", http.StatusBadRequest)
				return
			}
			priceID := fullSub.Items.Data[0].Price.ID
			fmt.Println("Price ID:", priceID)

			orgID := evtSub.Metadata["org_id"]
			if orgID == "" {
				http.Error(w, "Missing org ID in subscription metadata", http.StatusBadRequest)
				return
			}

			orgUUID, err := uuid.Parse(orgID)
			if err != nil {
				http.Error(w, "Error parsing ID to UUID", http.StatusInternalServerError)
				return
			}

			plan, ok := cfg.AllowedPlans[priceID]
			if !ok {
				plan = "basic" // fallback if priceID not found
			}

			if evtSub.Status == "active" {
				err = store.UpdateOrganizationBilling(r.Context(), db.UpdateOrganizationBillingParams{
					ID: orgUUID,
					StripeCustomerID: sql.NullString{
						String: evtSub.Customer.ID,
						Valid:  true,
					},
					StripeSubscriptionID: sql.NullString{
						String: evtSub.ID,
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

			org, err := store.GetOrganization(r.Context(), orgUUID)
			if err != nil {
				http.Error(w, "Error fetching Organization", http.StatusBadRequest)
				return
			}

			user, err := store.GetUserByID(r.Context(), org.OwnerID)
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				return
			}

			go func() {
				data := struct {
					Email string
					Plan  string
				}{
					Email: user.Email,
					Plan:  plan,
				}

				err := sender.SendTemplate(user.Email, "Your SaaSy Plan Changed", "plan_change.html.tmpl", data)
				if err != nil {
					log.Printf("Failed to send plan change email: %v", err)
				}
				fmt.Println("Email sent successfully")
			}()

			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusOK)
		}
	}
}
