package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"app-server/pkg/handlers"
	"app-server/pkg/s3"
	"app-server/pkg/templates"
)

func main() {
	// 获取端口配置，默认为8080（App Runner标准）
	port := os.Getenv("PORT_APP_SERVER") // Use PORT_APP_SERVER
	if port == "" {
		port = "8080"
	}

	fmt.Printf("📱 App Server Starting on :%s\n", port)

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
	appHandler := handlers.NewAppHandler(renderer, s3Client)

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
	fmt.Printf("🌐 Visit: http://localhost:%s\n", port)

	// Start server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
