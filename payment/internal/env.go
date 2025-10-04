package internal

import (
	"os"
)

const (
	PAYMENT_PORT = "PAYMENT_PORT"
	ORDER_URL	= "ORDER_URL"

	STRIPE_API_KEY = "STRIPE_API_KEY"
)

var (
	PaymentPort string
	OrderURL string
	StripeApiKey string
)

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(key + " environment variable is not set")
	}
	return val
}

func init() {
	// initialize environment variables
	PaymentPort = getEnv(PAYMENT_PORT)
	OrderURL = getEnv(ORDER_URL)
	StripeApiKey = getEnv(STRIPE_API_KEY)
}
