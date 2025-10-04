package internal

import (
	"os"
)

const (
	CATALOG_POSTGRES_URL = "CATALOG_POSTGRES_URL"
	CATALOG_PORT = "CATALOG_PORT"
)

var (
	CatalogPostgresURL string
	CatalogPort string
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
	CatalogPostgresURL = getEnv(CATALOG_POSTGRES_URL)
	CatalogPort = getEnv(CATALOG_PORT)
}
