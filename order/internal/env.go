package internal

import (
	"os"
)

const (
	ORDER_POSTGRES_URL = "ORDER_POSTGRES_URL"
	ORDER_PORT         = "ORDER_PORT"

	INVENTORY_URL = "INVENTORY_URL"
	CATALOG_URL   = "CATALOG_URL"
)

var (
	OrderPostgresURL string
	OrderPort         string
	InventoryURL string
	CatalogURL string
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
	OrderPostgresURL = getEnv(ORDER_POSTGRES_URL)
	OrderPort = getEnv(ORDER_PORT)
	InventoryURL = getEnv(INVENTORY_URL)
	CatalogURL = getEnv(CATALOG_URL)
}
