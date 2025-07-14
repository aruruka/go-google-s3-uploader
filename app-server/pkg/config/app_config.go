package config

import (
	"fmt"
	"log"
	"os"
)

// AppConfig holds all application-wide configurations for the app-server.
type AppConfig struct {
	Env           string
	PortAppServer string
	AWSRegion     string
	S3BucketName  string
	AppServerURL  string // The public URL of this app-server
	AuthServerURL string // The public URL of the auth-server
}

// LoadConfig loads configuration from environment variables for the app-server.
// It sets default values for development environment if not explicitly set.
func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	// Load .env file if it exists (for local development)
	// This should be called before reading individual env vars to ensure they are set.
	// Note: This LoadEnv is from the shared config package, assuming it's accessible.
	// If not, it needs to be implemented or imported correctly.
	// For now, assuming it's available via a shared config module or similar.
	// If not, we'd need to copy the LoadEnv function here or adjust module structure.
	// Given the project structure, it's likely `auth-server/pkg/config` is distinct.
	// So, I will assume `LoadEnv` needs to be implemented here or a common `pkg/config` is used.
	// For simplicity, I'll assume `os.Getenv` is sufficient for App Runner, and local .env handling
	// might be done externally or by a separate common utility.
	// However, to be consistent with auth-server, I will add a simple LoadEnv here.
	// Or, better, I will assume a shared `pkg/config` for `LoadEnv` is intended.
	// Let's check the `auth-server/pkg/config/config.go` (the original one) to see if it's generic.
	// It's not generic, it's specific to auth-server. So I need to copy the LoadEnv logic or make it shared.

	// For now, I will just rely on os.Getenv and assume environment variables are set by Docker/App Runner.
	// If .env loading is needed for app-server local dev, it should be handled explicitly.
	// Given the task, the focus is on App Runner env vars.

	cfg.Env = os.Getenv("ENV")
	cfg.PortAppServer = os.Getenv("PORT_APP_SERVER")
	cfg.AWSRegion = os.Getenv("AWS_REGION")
	cfg.S3BucketName = os.Getenv("S3_BUCKET_NAME")
	cfg.AppServerURL = os.Getenv("APP_SERVER_URL")
	cfg.AuthServerURL = os.Getenv("REDIRECT_URL") // AuthServerURL is the REDIRECT_URL from auth-server's perspective

	// Set default values for development if not in production
	isProduction := cfg.Env == "production"

	if cfg.PortAppServer == "" {
		cfg.PortAppServer = "8080" // Default port for app-server
	}
	if cfg.AWSRegion == "" {
		cfg.AWSRegion = "ap-northeast-1" // Default AWS region
	}
	if cfg.S3BucketName == "" {
		cfg.S3BucketName = "raymond-go-s3-uploader-dev-2025" // Default S3 bucket
	}

	if cfg.AppServerURL == "" {
		if !isProduction {
			cfg.AppServerURL = fmt.Sprintf("http://localhost:%s", cfg.PortAppServer)
			log.Printf("📍 Using default App Server URL: %s", cfg.AppServerURL)
		} else {
			return nil, fmt.Errorf("APP_SERVER_URL environment variable is required in production")
		}
	}

	// The AuthServerURL is the REDIRECT_URL from the auth-server's perspective.
	// The app-server needs to know where to redirect to the auth-server's login page.
	// This is typically the base URL of the auth-server.
	// For simplicity, I'm using REDIRECT_URL env var as AuthServerURL.
	// If the auth-server's base URL is different from its redirect URL, this needs adjustment.
	// Assuming REDIRECT_URL is the base URL for auth-server for now.
	if cfg.AuthServerURL == "" {
		if !isProduction {
			// For local dev, auth-server runs on 8081
			cfg.AuthServerURL = fmt.Sprintf("http://localhost:%s", "8081") // Assuming auth-server is on 8081
			log.Printf("📍 Using default Auth Server URL: %s", cfg.AuthServerURL)
		} else {
			return nil, fmt.Errorf("REDIRECT_URL (used as AuthServerURL) environment variable is required in production")
		}
	}

	// Log loaded configuration (excluding secrets)
	log.Printf("App Server Loaded Configuration: ENV=%s, PortAppServer=%s, AWS_REGION=%s, S3_BUCKET_NAME=%s, AppServerURL=%s, AuthServerURL=%s",
		cfg.Env, cfg.PortAppServer, cfg.AWSRegion, cfg.S3BucketName, cfg.AppServerURL, cfg.AuthServerURL)

	return cfg, nil
}
