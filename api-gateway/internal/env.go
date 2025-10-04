package internal

import (
	"os"
)

const (
	PROTOCOL = "PROTOCOL"
	PORT = "PORT"
	AUTH_URI = "AUTH_URI"
	CATALOG_URI = "CATALOG_URI"
	CART_URI = "CART_URI"
	INVENTORY_URI = "INVENTORY_URI"
	ORDER_URI = "ORDER_URI"
	PAYMENT_URI = "PAYMENT_URI"

	ROOT_PAGE_URI = "ROOT_PAGE_URI"
	LOGIN_URI = "LOGIN_URI"
)

var (
	Protocol string
	Port string
	AuthURI string
	CatalogURI string
	CartURI string
	InventoryURI string
	OrderURI string
	PaymentURI string

	RootPageURI string
	LoginURI string
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
	Protocol = getEnv(PROTOCOL)
	Port = getEnv(PORT)

	AuthURI = getEnv(AUTH_URI)
	CatalogURI = getEnv(CATALOG_URI)
	CartURI = getEnv(CART_URI)
	InventoryURI = getEnv(INVENTORY_URI)
	OrderURI = getEnv(ORDER_URI)
	PaymentURI = getEnv(PAYMENT_URI)

	RootPageURI = getEnv(ROOT_PAGE_URI)
	LoginURI = getEnv(LOGIN_URI)
}
