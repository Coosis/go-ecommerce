package internal

import (
	"os"
)

const (
	INVENTORY_POSTGRES_URL = "INVENTORY_POSTGRES_URL"
	INVENTORY_PORT         = "INVENTORY_PORT"
)

var (
	InventoryPostgresURL string
	InventoryPort         string
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
	InventoryPostgresURL = getEnv(INVENTORY_POSTGRES_URL)
	InventoryPort = getEnv(INVENTORY_PORT)
}
