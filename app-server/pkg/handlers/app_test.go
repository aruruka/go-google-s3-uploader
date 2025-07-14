package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"app-server/pkg/config" // Import the config package
)

// MockS3Client for testing handlers
type MockS3Client struct {
	UploadFileFunc    func(ctx context.Context, key string, file io.Reader, contentType string) error
	GetFileURLFunc    func(key string) string
	DeleteFileFunc    func(ctx context.Context, key string) error
	ListFilesFunc     func(ctx context.Context, prefix string) ([]string, error)
	ShouldReturnError bool
}

func (m *MockS3Client) UploadFile(ctx context.Context, key string, file io.Reader, contentType string) error {
	if m.ShouldReturnError {
		return errors.New("mock S3 upload error")
	}
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(ctx, key, file, contentType)
	}
	return nil
}

func (m *MockS3Client) GetFileURL(key string) string {
	if m.GetFileURLFunc != nil {
		return m.GetFileURLFunc(key)
	}
	return "https://mock-bucket.s3.amazonaws.com/" + key
}

func (m *MockS3Client) DeleteFile(ctx context.Context, key string) error {
	if m.ShouldReturnError {
		return errors.New("mock S3 delete error")
	}
	if m.DeleteFileFunc != nil {
		return m.DeleteFileFunc(ctx, key)
	}
	return nil
}

func (m *MockS3Client) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	if m.ShouldReturnError {
		return nil, errors.New("mock S3 list error")
	}
	if m.ListFilesFunc != nil {
		return m.ListFilesFunc(ctx, prefix)
	}
	return []string{}, nil
}

// MockTemplateRenderer for testing
type MockTemplateRenderer struct {
	ShouldReturnError bool
}

func (m *MockTemplateRenderer) RenderTemplate(w io.Writer, templateName string, data interface{}) error {
	if m.ShouldReturnError {
		return errors.New("mock template render error")
	}
	_, err := w.Write([]byte("mock template output"))
	return err
}

// Test AppHandler creation
func TestNewAppHandler(t *testing.T) {
	mockRenderer := &MockTemplateRenderer{}
	mockS3Client := &MockS3Client{}
	mockAppConfig := &config.AppConfig{
		AuthServerURL: "http://mock-auth-server.com", // Mock URL for testing redirects
		S3BucketName:  "mock-s3-bucket",
	}

	handler := NewAppHandler(mockAppConfig, mockRenderer, mockS3Client)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

// Test file type validation
func TestAppHandler_IsValidFileType(t *testing.T) {
	mockRenderer := &MockTemplateRenderer{}
	mockS3Client := &MockS3Client{}
	mockAppConfig := &config.AppConfig{
		AuthServerURL: "http://mock-auth-server.com",
		S3BucketName:  "mock-s3-bucket",
	}
	handler := &AppHandler{
		appConfig: mockAppConfig, // Add appConfig
		renderer:  mockRenderer,
		s3Client:  mockS3Client,
	}

	tests := []struct {
		contentType string
		expected    bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", true},
		{"image/webp", true},
		{"application/pdf", true},
		{"application/zip", true},
		{"application/x-zip-compressed", true},
		{"text/plain", false},
		{"video/mp4", false},
		{"application/exe", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			result := handler.isValidFileType(tt.contentType)
			if result != tt.expected {
				t.Errorf("isValidFileType(%s) = %v, expected %v", tt.contentType, result, tt.expected)
			}
		})
	}
}

// Test HandleHome
func TestAppHandler_HandleHome(t *testing.T) {
	mockRenderer := &MockTemplateRenderer{}
	mockS3Client := &MockS3Client{}
	mockAppConfig := &config.AppConfig{
		AuthServerURL: "http://mock-auth-server.com", // Mock URL for testing redirects
		S3BucketName:  "mock-s3-bucket",
	}
	handler := &AppHandler{
		appConfig: mockAppConfig, // Add appConfig
		renderer:  mockRenderer,
		s3Client:  mockS3Client,
	}

	req := httptest.NewRequest("GET", "/", nil)
	// Add session cookie for authenticated user - use "user_session" as cookie name
	req.AddCookie(&http.Cookie{
		Name:  "user_session",
		Value: "eyJpZCI6InRlc3QtdXNlci1pZCIsIm5hbWUiOiJKb2huIERvZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9", // base64 encoded JSON
	})

	w := httptest.NewRecorder()
	handler.HandleHome(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
