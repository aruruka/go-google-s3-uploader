package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/config" // Import the new config package
	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/handlers"
	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/s3"
	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/templates"
)

func main() {
	fmt.Println("📱 App Server Starting...")

	// Load application configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application configuration: %v", err)
	}

	fmt.Printf("📱 App Server Starting on :%s\n", appConfig.PortAppServer)

	// Initialize template renderer
	renderer, err := templates.NewTemplateRenderer()
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}

	// Initialize S3 client
	s3Client, err := s3.NewS3Client()
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Initialize handlers
	appHandler := handlers.NewAppHandler(appConfig, renderer, s3Client) // Pass appConfig

	// Define routes
	http.HandleFunc("/", appHandler.HandleHome)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			appHandler.HandleUpload(w, r)
		case http.MethodPost:
			appHandler.HandleUploadPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// API endpoint for file upload (used by frontend form)
	http.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			appHandler.HandleUploadPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/success", appHandler.HandleSuccess)

	// 健康检查端点 (App Runner 要求)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../shared/static/"))))

	fmt.Println("✅ App Server ready")
	fmt.Printf("🌐 Visit: %s\n", appConfig.AppServerURL) // Use AppServerURL for visit message

	// Start server
	log.Fatal(http.ListenAndServe(":"+appConfig.PortAppServer, nil))
}
