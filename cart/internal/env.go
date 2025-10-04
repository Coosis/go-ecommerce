package internal

import (
	"os"
)

const (
	CART_POSTGRES_URL = "CART_POSTGRES_URL"
	CART_PORT = "CART_PORT"
	CATALOG_URL = "CATALOG_URL"
	ORDER_URL = "ORDER_URL"
)

var (
	OrderURL string
	CatalogURL string
	CartPort string
	CartPostgresURL string
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
	OrderURL = getEnv(OrderURL)
	CatalogURL = getEnv(CatalogURL)
	CartPort = getEnv(CartPort)
	CartPostgresURL = getEnv(CartPostgresURL)
}
