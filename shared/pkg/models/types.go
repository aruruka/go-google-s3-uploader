package models

import "time"

// User represents a user in the system
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Picture  string    `json:"picture"`
	Provider string    `json:"provider"`
	Created  time.Time `json:"created"`
}

// PageData represents the common data structure for all pages
type PageData struct {
	Title        string      `json:"title"`
	User         *User       `json:"user,omitempty"`
	FlashMessage string      `json:"flash_message,omitempty"`
	FlashType    string      `json:"flash_type,omitempty"` // success, error, info, warning
	Data         interface{} `json:"data,omitempty"`       // Page-specific data
	CSRFToken    string      `json:"csrf_token,omitempty"`
}

// ErrorData represents data for error pages
type ErrorData struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

// FileUpload represents an uploaded file
type FileUpload struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	S3Key       string    `json:"s3_key"`
	S3URL       string    `json:"s3_url"`
	UploadedAt  time.Time `json:"uploaded_at"`
	UserID      string    `json:"user_id"`
}

// UploadResponse represents the response after file upload
type UploadResponse struct {
	Success bool        `json:"success"`
	File    *FileUpload `json:"file,omitempty"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
}

// HomeData represents data for the home page
type HomeData struct {
	RecentUploads []FileUpload `json:"recent_uploads"`
	TotalUploads  int          `json:"total_uploads"`
	TotalSize     int64        `json:"total_size"`
}

// UploadData represents data for the upload page
type UploadData struct {
	MaxFileSize  int64    `json:"max_file_size"`
	AllowedTypes []string `json:"allowed_types"`
	S3BucketName string   `json:"s3_bucket_name,omitempty"`
}

// SuccessData represents data for the success page
type SuccessData struct {
	Upload      *FileUpload `json:"upload"`
	RedirectURL string      `json:"redirect_url,omitempty"`
}

// Auth-specific models

// LoginData represents data for the login page
type LoginData struct {
	RedirectURL string `json:"redirect_url,omitempty"`
	Error       string `json:"error,omitempty"`
}

// CallbackData represents data for the OAuth callback page
type CallbackData struct {
	User        *User  `json:"user,omitempty"`
	Error       string `json:"error,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}
