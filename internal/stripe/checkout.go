package stripeclient

import (
	"fmt"
	"os"

	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CreateCheckoutSeassion(orgID string, priceID string) (string, error) {
	SuccessURL := fmt.Sprintf("%s/success", os.Getenv("FRONTEND_URL"))
	CancelURL := fmt.Sprintf("%s/cancel", os.Getenv("FRONTEND_URL"))

	params := &stripe.CheckoutSessionParams{
		SuccessURL:        stripe.String(SuccessURL),
		CancelURL:         stripe.String(CancelURL),
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		ClientReferenceID: stripe.String(orgID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"org_id": orgID,
			},
		},
	}

	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}
