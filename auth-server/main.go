package main

import (
	"fmt"
	"log"
	"net/http"

	"auth-server/pkg/config" // This now refers to the package containing LoadEnv and AppConfig
	"auth-server/pkg/handlers"
	"auth-server/pkg/oauth"
	"auth-server/pkg/templates"
)

func main() {
	fmt.Println("🔐 Auth Server Starting...")

	// Load application configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application configuration: %v", err)
	}

	// Initialize OAuth configuration
	oauthConfig, err := oauth.NewConfig(appConfig)
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
	fmt.Printf("🌐 Visit: http://localhost:%s/login\n", appConfig.PortAuthServer)

	// Start server
	log.Fatal(http.ListenAndServe(":"+appConfig.PortAuthServer, mux))
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	http.NotFound(w, r)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"healthy","service":"auth-server"}`)
}
