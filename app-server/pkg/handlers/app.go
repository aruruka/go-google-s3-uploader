package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/config"
	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/s3"
	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/templates"
	"github.com/aruruka/go-google-s3-uploader/shared/pkg/models"
)

// AppHandlerIface defines the interface for application handlers
type AppHandlerIface interface {
	HandleHome(w http.ResponseWriter, r *http.Request)
	HandleUpload(w http.ResponseWriter, r *http.Request)
	HandleUploadPost(w http.ResponseWriter, r *http.Request)
	HandleSuccess(w http.ResponseWriter, r *http.Request)
}

// AppHandler implements application handlers
type AppHandler struct {
	appConfig *config.AppConfig // Add appConfig
	renderer  templates.TemplateRendererIface
	s3Client  s3.S3ClientIface
}

// NewAppHandler creates a new application handler
func NewAppHandler(appConfig *config.AppConfig, renderer templates.TemplateRendererIface, s3Client s3.S3ClientIface) AppHandlerIface {
	return &AppHandler{
		appConfig: appConfig, // Store appConfig
		renderer:  renderer,
		s3Client:  s3Client,
	}
}

// HandleHome displays the home page
func (h *AppHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("🏠 App-server HandleHome called from: %s", r.RemoteAddr)

	// Check if user is authenticated
	user := h.getUserFromSession(r)
	if user == nil {
		log.Printf("❌ No user session found, redirecting to auth-server")
		// Check if AuthServerURL is the same as our domain (App Runner scenario)
		if h.appConfig.AuthServerURL == h.appConfig.AppServerURL {
			// Internal redirect to login page within the same service
			log.Printf("🔄 Internal redirect to /login (same domain)")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		} else {
			// External redirect to separate auth server (local development)
			http.Redirect(w, r, h.appConfig.AuthServerURL+"/login", http.StatusTemporaryRedirect)
		}
		return
	}

	log.Printf("✅ User session found: %s (%s)", user.Name, user.Email)

	// Prepare page data
	pageData := &models.PageData{
		Title: "Google S3 Uploader - Home",
		User:  user,
		Data: &models.HomeData{
			RecentUploads: []models.FileUpload{}, // TODO: Load from database
			TotalUploads:  0,
			TotalSize:     0,
			AuthServerURL: h.appConfig.AuthServerURL, // Pass AuthServerURL to template
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.RenderTemplate(w, "home.html", pageData); err != nil {
		log.Printf("Failed to render home template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleUpload displays the upload page
func (h *AppHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	user := h.getUserFromSession(r)
	if user == nil {
		// Redirect to auth server for login
		http.Redirect(w, r, h.appConfig.AuthServerURL+"/login", http.StatusTemporaryRedirect) // Use appConfig
		return
	}

	pageData := &models.PageData{
		Title: "Upload File - Google S3 Uploader",
		User:  user,
		Data: &models.UploadData{
			MaxFileSize:  50 * 1024 * 1024, // 50 MB
			AllowedTypes: []string{"image/*", "application/pdf", "application/zip"},
			S3BucketName: h.appConfig.S3BucketName, // Use appConfig
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.RenderTemplate(w, "upload.html", pageData); err != nil {
		log.Printf("Failed to render upload template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleUploadPost processes file upload
func (h *AppHandler) HandleUploadPost(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	user := h.getUserFromSession(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(50 << 20) // 50 MB max
	if err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		h.renderError(w, "Failed to parse upload form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
		h.renderError(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO: Validate file type and size
	if fileHeader.Size > 50*1024*1024 { // 50 MB
		h.renderError(w, "File too large (max 50 MB)", http.StatusBadRequest)
		return
	}

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	if !h.isValidFileType(contentType) {
		h.renderError(w, "Invalid file type. Only images, PDFs, and ZIP files are allowed", http.StatusBadRequest)
		return
	}

	// Generate S3 key
	s3Key := fmt.Sprintf("uploads/%s/%d_%s", user.ID, time.Now().Unix(), fileHeader.Filename)

	// Upload to S3
	ctx := context.Background()
	err = h.s3Client.UploadFile(ctx, s3Key, file, contentType)
	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		h.renderError(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// Create upload record
	uploadedFile := &models.FileUpload{
		ID:          fmt.Sprintf("file_%d", time.Now().Unix()),
		Filename:    fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: contentType,
		S3Key:       s3Key,
		S3URL:       h.s3Client.GetFileURL(s3Key),
		UploadedAt:  time.Now(),
		UserID:      user.ID,
	}

	log.Printf("File uploaded successfully: %s (%d bytes)", uploadedFile.Filename, uploadedFile.Size)

	// For now, we'll pass the file info via query parameters
	// In production, this would be stored in a database
	queryParams := fmt.Sprintf("?filename=%s&size=%d&contentType=%s&s3url=%s&uploadTime=%s",
		uploadedFile.Filename,
		uploadedFile.Size,
		uploadedFile.ContentType,
		uploadedFile.S3URL,
		uploadedFile.UploadedAt.Format("2006-01-02 15:04:05"),
	)

	// Redirect to success page
	http.Redirect(w, r, "/success"+queryParams, http.StatusSeeOther)
}

// HandleSuccess displays the success page
func (h *AppHandler) HandleSuccess(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	user := h.getUserFromSession(r)
	if user == nil {
		http.Redirect(w, r, h.appConfig.AuthServerURL+"/login", http.StatusTemporaryRedirect) // Use appConfig
		return
	}

	// Get file details from query parameters
	filename := r.URL.Query().Get("filename")
	sizeStr := r.URL.Query().Get("size")
	contentType := r.URL.Query().Get("contentType")
	s3url := r.URL.Query().Get("s3url")
	uploadTimeStr := r.URL.Query().Get("uploadTime")

	if filename == "" || sizeStr == "" {
		h.renderError(w, "File information not found", http.StatusNotFound)
		return
	}

	// Parse size
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		log.Printf("Failed to parse file size: %v", err)
		size = 0
	}

	// Parse upload time
	uploadTime, err := time.Parse("2006-01-02 15:04:05", uploadTimeStr)
	if err != nil {
		log.Printf("Failed to parse upload time: %v", err)
		uploadTime = time.Now()
	}

	// Create file upload data
	uploadedFile := &models.FileUpload{
		ID:          fmt.Sprintf("file_%d", uploadTime.Unix()),
		Filename:    filename,
		Size:        size,
		ContentType: contentType,
		S3URL:       s3url,
		UploadedAt:  uploadTime,
		UserID:      user.ID,
	}

	pageData := &models.PageData{
		Title: "Upload Successful - Google S3 Uploader",
		User:  user,
		Data: &models.SuccessData{
			Upload:      uploadedFile,
			RedirectURL: "/",
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.RenderTemplate(w, "success.html", pageData); err != nil {
		log.Printf("Failed to render success template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// getUserFromSession extracts user from session cookie
func (h *AppHandler) getUserFromSession(r *http.Request) *models.User {
	log.Printf("🍪 Checking for user_session cookie...")

	// Debug: print all cookies
	for _, cookie := range r.Cookies() {
		log.Printf("🍪 Found cookie: %s = %s", cookie.Name, cookie.Value[:min(50, len(cookie.Value))])
	}

	cookie, err := r.Cookie("user_session")
	if err != nil {
		log.Printf("❌ user_session cookie not found: %v", err)
		return nil
	}

	log.Printf("✅ user_session cookie found, decoding...")

	// Decode base64 session data
	sessionData, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Printf("Failed to decode session data: %v", err)
		return nil
	}

	// Parse user JSON
	var user models.User
	if err := json.Unmarshal(sessionData, &user); err != nil {
		log.Printf("Failed to parse user session: %v", err)
		return nil
	}

	log.Printf("✅ Session decoded successfully: %s", user.Name)
	return &user
}

// isValidFileType checks if the content type is allowed
func (h *AppHandler) isValidFileType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":                   true,
		"image/png":                    true,
		"image/gif":                    true,
		"image/webp":                   true,
		"application/pdf":              true,
		"application/zip":              true,
		"application/x-zip-compressed": true,
	}
	return allowedTypes[contentType]
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// renderError renders an error page
func (h *AppHandler) renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)

	pageData := &models.PageData{
		Title: "Error",
		Data: &models.ErrorData{
			StatusCode: statusCode,
			Message:    message,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.RenderTemplate(w, "error.html", pageData); err != nil {
		log.Printf("Failed to render error template: %v", err)
		// Fallback to plain text error
		http.Error(w, message, statusCode)
	}
}
