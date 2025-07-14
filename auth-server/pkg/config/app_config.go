package config

import (
	"fmt"
	"log"
	"os"
)

// AppConfig holds all application-wide configurations.
type AppConfig struct {
	Env                string
	PortAuthServer     string // Port for auth-server
	PortAppServer      string // Port for app-server
	AWSRegion          string
	S3BucketName       string
	GoogleClientID     string
	GoogleClientSecret string
	RedirectURL        string
	AppServerURL       string
}

// LoadConfig loads configuration from environment variables.
// It sets default values for development environment if not explicitly set.
func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	// Load .env file if it exists (for local development)
	// This should be called before reading individual env vars to ensure they are set.
	if err := LoadEnv(".env"); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	cfg.Env = os.Getenv("ENV")
	cfg.PortAuthServer = os.Getenv("PORT_AUTH_SERVER")
	cfg.PortAppServer = os.Getenv("PORT_APP_SERVER")
	cfg.AWSRegion = os.Getenv("AWS_REGION")
	cfg.S3BucketName = os.Getenv("S3_BUCKET_NAME")
	cfg.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.RedirectURL = os.Getenv("REDIRECT_URL")
	cfg.AppServerURL = os.Getenv("APP_SERVER_URL")

	// Set default values for development if not in production
	isProduction := cfg.Env == "production"

	if cfg.PortAuthServer == "" {
		cfg.PortAuthServer = "8081" // Default port for auth-server
	}
	if cfg.PortAppServer == "" {
		cfg.PortAppServer = "8080" // Default port for app-server
	}
	if cfg.AWSRegion == "" {
		cfg.AWSRegion = "ap-northeast-1" // Default AWS region
	}
	if cfg.S3BucketName == "" {
		cfg.S3BucketName = "raymond-go-s3-uploader-dev-2025" // Default S3 bucket
	}

	if cfg.GoogleClientID == "" {
		log.Println("⚠️  GOOGLE_CLIENT_ID not set.")
		if !isProduction {
			cfg.GoogleClientID = "your-google-client-id" // Default for dev
		} else {
			return nil, fmt.Errorf("GOOGLE_CLIENT_ID environment variable is required in production")
		}
	}
	if cfg.GoogleClientSecret == "" {
		log.Println("⚠️  GOOGLE_CLIENT_SECRET not set.")
		if !isProduction {
			cfg.GoogleClientSecret = "your-google-client-secret" // Default for dev
		} else {
			return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET environment variable is required in production")
		}
	}
	if cfg.RedirectURL == "" {
		if !isProduction {
			cfg.RedirectURL = fmt.Sprintf("http://localhost:%s/auth/callback", cfg.PortAuthServer)
			log.Printf("📍 Using default redirect URL: %s", cfg.RedirectURL)
		} else {
			return nil, fmt.Errorf("REDIRECT_URL environment variable is required in production")
		}
	}
	if cfg.AppServerURL == "" {
		if !isProduction {
			cfg.AppServerURL = fmt.Sprintf("http://localhost:%s", cfg.PortAppServer) // Default for dev
			log.Printf("📍 Using default App Server URL: %s", cfg.AppServerURL)
		} else {
			return nil, fmt.Errorf("APP_SERVER_URL environment variable is required in production")
		}
	}

	// Log loaded configuration (excluding secrets)
	log.Printf("Loaded Configuration: ENV=%s, PortAuthServer=%s, PortAppServer=%s, AWS_REGION=%s, S3_BUCKET_NAME=%s, RedirectURL=%s, AppServerURL=%s",
		cfg.Env, cfg.PortAuthServer, cfg.PortAppServer, cfg.AWSRegion, cfg.S3BucketName, cfg.RedirectURL, cfg.AppServerURL)

	return cfg, nil
}
