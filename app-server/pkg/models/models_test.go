package models

import (
	"testing"
	"time"
)

// Test User model
func TestUser(t *testing.T) {
	user := User{
		ID:    "test-user-123",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if user.ID != "test-user-123" {
		t.Errorf("Expected ID 'test-user-123', got '%s'", user.ID)
	}
	
	if user.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", user.Name)
	}
	
	if user.Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got '%s'", user.Email)
	}
}

// Test FileUpload model
func TestFileUpload(t *testing.T) {
	now := time.Now()
	upload := FileUpload{
		ID:          "upload-123",
		Filename:    "test.jpg",
		Size:        1024,
		ContentType: "image/jpeg",
		S3Key:       "uploads/user123/test.jpg",
		S3URL:       "https://bucket.s3.amazonaws.com/uploads/user123/test.jpg",
		UploadedAt:  now,
		UserID:      "user-123",
	}

	if upload.ID != "upload-123" {
		t.Errorf("Expected ID 'upload-123', got '%s'", upload.ID)
	}
	
	if upload.Filename != "test.jpg" {
		t.Errorf("Expected Filename 'test.jpg', got '%s'", upload.Filename)
	}
	
	if upload.Size != 1024 {
		t.Errorf("Expected Size 1024, got %d", upload.Size)
	}
	
	if upload.ContentType != "image/jpeg" {
		t.Errorf("Expected ContentType 'image/jpeg', got '%s'", upload.ContentType)
	}
	
	if upload.UploadedAt != now {
		t.Errorf("Expected UploadedAt to match now")
	}
}

// Test PageData model
func TestPageData(t *testing.T) {
	user := &User{ID: "test", Name: "Test User", Email: "test@example.com"}
	pageData := PageData{
		Title: "Test Page",
		User:  user,
		Data:  map[string]interface{}{"key": "value"},
	}

	if pageData.Title != "Test Page" {
		t.Errorf("Expected Title 'Test Page', got '%s'", pageData.Title)
	}
	
	if pageData.User.Name != "Test User" {
		t.Errorf("Expected User Name 'Test User', got '%s'", pageData.User.Name)
	}
}

// Test different data types
func TestHomeData(t *testing.T) {
	homeData := HomeData{
		RecentUploads: []FileUpload{
			{ID: "1", Filename: "file1.jpg"},
			{ID: "2", Filename: "file2.png"},
		},
		TotalUploads: 2,
		TotalSize:    2048,
	}

	if len(homeData.RecentUploads) != 2 {
		t.Errorf("Expected 2 recent uploads, got %d", len(homeData.RecentUploads))
	}
	
	if homeData.TotalUploads != 2 {
		t.Errorf("Expected TotalUploads 2, got %d", homeData.TotalUploads)
	}
	
	if homeData.TotalSize != 2048 {
		t.Errorf("Expected TotalSize 2048, got %d", homeData.TotalSize)
	}
}

func TestUploadData(t *testing.T) {
	uploadData := UploadData{
		MaxFileSize:  50 * 1024 * 1024, // 50MB
		AllowedTypes: []string{"image/*", "application/pdf"},
		S3BucketName: "test-bucket",
	}

	if uploadData.MaxFileSize != 50*1024*1024 {
		t.Errorf("Expected MaxFileSize 50MB, got %d", uploadData.MaxFileSize)
	}
	
	if len(uploadData.AllowedTypes) != 2 {
		t.Errorf("Expected 2 allowed types, got %d", len(uploadData.AllowedTypes))
	}
	
	if uploadData.S3BucketName != "test-bucket" {
		t.Errorf("Expected S3BucketName 'test-bucket', got '%s'", uploadData.S3BucketName)
	}
}

func TestSuccessData(t *testing.T) {
	upload := &FileUpload{
		ID:       "test-upload",
		Filename: "success.jpg",
	}
	
	successData := SuccessData{
		Upload:      upload,
		RedirectURL: "/dashboard",
	}

	if successData.Upload.ID != "test-upload" {
		t.Errorf("Expected Upload ID 'test-upload', got '%s'", successData.Upload.ID)
	}
	
	if successData.RedirectURL != "/dashboard" {
		t.Errorf("Expected RedirectURL '/dashboard', got '%s'", successData.RedirectURL)
	}
}

func TestErrorData(t *testing.T) {
	errorData := ErrorData{
		StatusCode: 404,
		Message:    "File not found",
	}

	if errorData.StatusCode != 404 {
		t.Errorf("Expected StatusCode 404, got %d", errorData.StatusCode)
	}
	
	if errorData.Message != "File not found" {
		t.Errorf("Expected Message 'File not found', got '%s'", errorData.Message)
	}
}
