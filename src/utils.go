package main

import (
	"os"
)

// getPort returns the port by first looking at any environment variable
// nammed PORT and then defaulting to :8000
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8000"
}
