package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"auth-server/pkg/config"
	"auth-server/pkg/handlers"
	"auth-server/pkg/oauth"
	"auth-server/pkg/templates"
)

func main() {
	fmt.Println("🔐 Auth Server Starting on :8081")

	// Load environment variables from .env file
	if err := config.LoadEnv(".env"); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	// Set default environment variables for development if not already set
	setDefaultEnvVars()

	// Initialize OAuth configuration
	oauthConfig, err := oauth.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize OAuth config: %v", err)
	}

	// Initialize template renderer
	renderer, err := templates.NewTemplateRenderer()
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	// Initialize handlers with dependency injection
	authHandler := handlers.NewAuthHandler(oauthConfig, renderer)

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/", redirectToLogin)
	mux.HandleFunc("/login", authHandler.HandleLogin)
	mux.HandleFunc("/auth/google", authHandler.HandleGoogleAuth)
	mux.HandleFunc("/auth/callback", authHandler.HandleCallback)
	mux.HandleFunc("/logout", authHandler.HandleLogout)

	// Health check
	mux.HandleFunc("/health", healthCheck)

	// Serve static files
	fs := http.FileServer(http.Dir("../shared/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("✅ Auth Server ready")
	fmt.Println("🌐 Visit: http://localhost:8081/login")

	// Start server
	log.Fatal(http.ListenAndServe(":8081", mux))
}

// setDefaultEnvVars sets default environment variables for development
func setDefaultEnvVars() {
	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		log.Println("⚠️  GOOGLE_CLIENT_ID not set - OAuth will not work")
		os.Setenv("GOOGLE_CLIENT_ID", "your-google-client-id")
	}
	if os.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		log.Println("⚠️  GOOGLE_CLIENT_SECRET not set - OAuth will not work")
		os.Setenv("GOOGLE_CLIENT_SECRET", "your-google-client-secret")
	}
	if os.Getenv("REDIRECT_URL") == "" {
		os.Setenv("REDIRECT_URL", "http://localhost:8081/auth/callback")
		log.Println("📍 Using default redirect URL: http://localhost:8081/auth/callback")
	}
}

// redirectToLogin redirects root requests to login page
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	http.NotFound(w, r)
}

// healthCheck provides a simple health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"healthy","service":"auth-server"}`)
}
