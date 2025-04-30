package billing

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/config"
	"github.com/bmkersey/Go-SaaSy/internal/orgs"
	stripeclient "github.com/bmkersey/Go-SaaSy/internal/stripe"
)

type CheckoutRequest struct {
	Plan string `json:"plan"`
}

type CheckoutResponse struct {
	URL string `json:"url"`
}

func CreateCheckoutHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		orgID, ok := orgs.GetOrgID(r)
		if !ok {
			http.Error(w, "You need an org to subscribe", http.StatusUnauthorized)
			return
		}

		var req CheckoutRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		priceID, ok := cfg.PlanPriceIDs[req.Plan]
		if !ok {
			http.Error(w, "Invalid plan", http.StatusBadRequest)
			return
		}
		log.Println("checkout req plan: " + req.Plan + "\ncheckout price id: " + priceID)

		sessionURL, err := stripeclient.CreateCheckoutSeassion(orgID, priceID)
		if err != nil {
			http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		}

		resp := CheckoutResponse{URL: sessionURL}
		json.NewEncoder(w).Encode(resp)

	}
}
