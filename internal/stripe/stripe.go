package stripeclient

import (
	"os"

	"github.com/stripe/stripe-go/v76"
)

func Init() {
	stripe.Key = os.Getenv("STRIPE_SECRET")
}
