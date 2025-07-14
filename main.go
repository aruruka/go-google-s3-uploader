package main

import (
	"log"
	"net/http"
	"os"

	// Auth server imports
	authConfig "github.com/aruruka/go-google-s3-uploader/auth-server/pkg/config"
	authHandlers "github.com/aruruka/go-google-s3-uploader/auth-server/pkg/handlers"
	authOAuth "github.com/aruruka/go-google-s3-uploader/auth-server/pkg/oauth"
	authTemplates "github.com/aruruka/go-google-s3-uploader/auth-server/pkg/templates"

	// App server imports
	appConfig "github.com/aruruka/go-google-s3-uploader/app-server/pkg/config"
	appHandlers "github.com/aruruka/go-google-s3-uploader/app-server/pkg/handlers"
	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/s3"
	appTemplates "github.com/aruruka/go-google-s3-uploader/app-server/pkg/templates"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	log.Println("🚀 Starting combined Go S3 Uploader service...")

	// Load auth server configuration
	authAppConfig, err := authConfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load auth server config: %v", err)
	}

	// Load app server configuration
	appAppConfig, err := appConfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load app server config: %v", err)
	}

	// Initialize auth server components
	authRenderer, err := authTemplates.NewTemplateRenderer()
	if err != nil {
		log.Fatalf("Failed to create auth renderer: %v", err)
	}
	oauthConfig, err := authOAuth.NewConfig(authAppConfig)
	if err != nil {
		log.Fatalf("Failed to create OAuth config: %v", err)
	}
	authHandler := authHandlers.NewAuthHandler(authAppConfig, oauthConfig, authRenderer)

	// Initialize app server components
	appRenderer, err := appTemplates.NewTemplateRenderer()
	if err != nil {
		log.Fatalf("Failed to create app renderer: %v", err)
	}
	s3Client, err := s3.NewS3Client()
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	appHandler := appHandlers.NewAppHandler(appAppConfig, appRenderer, s3Client)

	// Create combined router
	mux := http.NewServeMux()

	// Auth server routes
	mux.HandleFunc("/login", authHandler.HandleLogin)
	mux.HandleFunc("/auth/google", authHandler.HandleGoogleAuth)
	mux.HandleFunc("/auth/callback", authHandler.HandleCallback)
	mux.HandleFunc("/logout", authHandler.HandleLogout)

	// App server routes
	mux.HandleFunc("/", appHandler.HandleHome)
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			appHandler.HandleUpload(w, r)
		} else if r.Method == http.MethodPost {
			appHandler.HandleUploadPost(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			appHandler.HandleUploadPost(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/success", appHandler.HandleSuccess)

	// Shared routes
	mux.HandleFunc("/health", healthCheck)

	// Static file serving
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("shared/static/"))))

	// Determine port - App Runner sets PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		// Fallback to app server port for local development
		port = appAppConfig.PortAppServer
		if port == "" {
			port = "8080" // Default fallback
		}
	}

	log.Printf("🌐 Server starting on port %s", port)
	log.Printf("📍 Auth routes: /login, /auth/google, /auth/callback, /logout")
	log.Printf("📍 App routes: /, /upload, /api/upload, /success")
	log.Printf("🔧 Health check: /health")
	log.Printf("📁 Static files: /static/")

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
