package stripeclient

import (
	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CreateCheckoutSeassion(orgID string, priceID string) (string, error) {
	// SuccessURL := fmt.Sprintf("%s/success", os.Getenv("FRONTEND_URL"))
	// CancelURL := fmt.Sprintf("%s/cancel", os.Getenv("FRONTEND_URL"))

	params := &stripe.CheckoutSessionParams{
		SuccessURL:        stripe.String("https://example.com/success"),
		CancelURL:         stripe.String("https://example.com/cancel"),
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
	// params := &stripe.CheckoutSessionParams{
	// 	ClientReferenceID: stripe.String(orgID), // ✅ Important!
	// 	Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	// 	Customer:          stripe.String(orgID), // if you have it
	// 	SuccessURL:        stripe.String("https://example.com/success"),
	// 	CancelURL:         stripe.String("https://example.com/cancel"),
	// 	LineItems: []*stripe.CheckoutSessionLineItemParams{
	// 		{
	// 			Price:    stripe.String(priceID),
	// 			Quantity: stripe.Int64(1),
	// 		},
	// 	},
	// }

	params.AddExpand("line_items")

	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}
