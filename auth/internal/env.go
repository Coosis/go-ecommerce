package internal

import (
	"os"
)

const (
	VALKEY_URL = "VALKEY_URL"
	AUTH_PORT = "AUTH_PORT"
	AUTH_POSTGRES_URL = "AUTH_POSTGRES_URL"

	AMQP_URL = "AMQP_URL"
	AMQP_QUEUE = "AMQP_QUEUE"
)

var (
	ValkeyURL string
	AuthPort string
	AuthPostgresUrl string
	AmqpUrl string
	AmqpQueue string
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
	ValkeyURL = getEnv(VALKEY_URL)
	AuthPort = getEnv(AUTH_PORT)
	AuthPostgresUrl = getEnv(AUTH_POSTGRES_URL)
	AmqpUrl = getEnv(AMQP_URL)
	AmqpQueue = getEnv(AMQP_QUEUE)
}
