package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const DEFAULT_API_URL = "http://ec2-56-125-224-81.sa-east-1.compute.amazonaws.com"

// Config holds the application configuration
type Config struct {
	APIURL string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Get API URL from environment variable
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = DEFAULT_API_URL
	}

	return &Config{
		APIURL: apiURL,
	}, nil
}
