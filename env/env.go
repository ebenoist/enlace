package env

import (
	"log"
	"os"
)

func Get(key string, fallback ...string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	if len(fallback) == 0 {
		log.Fatalf("missing required environment variable: %s", key)
	}

	return fallback[0]
}
