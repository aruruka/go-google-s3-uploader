package s3

import (
	"context"
	"io"
	"strings"
	"testing"
)

// MockS3Client implements S3ClientIface for testing
type MockS3Client struct {
	uploadError   error
	deleteError   error
	listError     error
	uploadedFiles map[string][]byte
	baseURL       string
}

// NewMockS3Client creates a new mock S3 client
func NewMockS3Client() *MockS3Client {
	return &MockS3Client{
		uploadedFiles: make(map[string][]byte),
		baseURL:       "https://test-bucket.s3.ap-northeast-1.amazonaws.com",
	}
}

// UploadFile mocks file upload
func (m *MockS3Client) UploadFile(ctx context.Context, key string, file io.Reader, contentType string) error {
	if m.uploadError != nil {
		return m.uploadError
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	m.uploadedFiles[key] = content
	return nil
}

// GetFileURL mocks URL generation
func (m *MockS3Client) GetFileURL(key string) string {
	return m.baseURL + "/" + key
}

// DeleteFile mocks file deletion
func (m *MockS3Client) DeleteFile(ctx context.Context, key string) error {
	if m.deleteError != nil {
		return m.deleteError
	}

	delete(m.uploadedFiles, key)
	return nil
}

// ListFiles mocks file listing
func (m *MockS3Client) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	if m.listError != nil {
		return nil, m.listError
	}

	var files []string
	for key := range m.uploadedFiles {
		if strings.HasPrefix(key, prefix) {
			files = append(files, key)
		}
	}
	return files, nil
}

// Test helper methods
func (m *MockS3Client) SetUploadError(err error) {
	m.uploadError = err
}

func (m *MockS3Client) SetDeleteError(err error) {
	m.deleteError = err
}

func (m *MockS3Client) SetListError(err error) {
	m.listError = err
}

func (m *MockS3Client) GetUploadedContent(key string) []byte {
	return m.uploadedFiles[key]
}

func (m *MockS3Client) IsFileUploaded(key string) bool {
	_, exists := m.uploadedFiles[key]
	return exists
}

// Unit tests for S3Client
func TestMockS3Client_UploadFile(t *testing.T) {
	mockClient := NewMockS3Client()

	tests := []struct {
		name        string
		key         string
		content     string
		contentType string
		wantError   bool
	}{
		{
			name:        "successful upload",
			key:         "test/file.txt",
			content:     "Hello, World!",
			contentType: "text/plain",
			wantError:   false,
		},
		{
			name:        "upload with special characters",
			key:         "uploads/用户123/文件.pdf",
			content:     "PDF content",
			contentType: "application/pdf",
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.content)
			err := mockClient.UploadFile(context.Background(), tt.key, reader, tt.contentType)

			if (err != nil) != tt.wantError {
				t.Errorf("UploadFile() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if !mockClient.IsFileUploaded(tt.key) {
					t.Errorf("File was not uploaded: %s", tt.key)
				}

				uploadedContent := mockClient.GetUploadedContent(tt.key)
				if string(uploadedContent) != tt.content {
					t.Errorf("Uploaded content mismatch. Got: %s, Want: %s", string(uploadedContent), tt.content)
				}
			}
		})
	}
}

func TestMockS3Client_GetFileURL(t *testing.T) {
	mockClient := NewMockS3Client()

	key := "uploads/user123/test.jpg"
	expectedURL := "https://test-bucket.s3.ap-northeast-1.amazonaws.com/uploads/user123/test.jpg"

	actualURL := mockClient.GetFileURL(key)
	if actualURL != expectedURL {
		t.Errorf("GetFileURL() = %v, want %v", actualURL, expectedURL)
	}
}

func TestMockS3Client_DeleteFile(t *testing.T) {
	mockClient := NewMockS3Client()

	// First upload a file
	key := "test/delete-me.txt"
	content := strings.NewReader("test content")
	err := mockClient.UploadFile(context.Background(), key, content, "text/plain")
	if err != nil {
		t.Fatalf("Failed to upload test file: %v", err)
	}

	// Verify file exists
	if !mockClient.IsFileUploaded(key) {
		t.Fatal("Test file was not uploaded")
	}

	// Delete the file
	err = mockClient.DeleteFile(context.Background(), key)
	if err != nil {
		t.Errorf("DeleteFile() error = %v", err)
	}

	// Verify file is deleted
	if mockClient.IsFileUploaded(key) {
		t.Error("File was not deleted")
	}
}

func TestMockS3Client_ListFiles(t *testing.T) {
	mockClient := NewMockS3Client()

	// Upload test files
	testFiles := []string{
		"uploads/user1/file1.txt",
		"uploads/user1/file2.pdf",
		"uploads/user2/file3.jpg",
		"temp/file4.zip",
	}

	for _, key := range testFiles {
		content := strings.NewReader("test content")
		err := mockClient.UploadFile(context.Background(), key, content, "application/octet-stream")
		if err != nil {
			t.Fatalf("Failed to upload test file %s: %v", key, err)
		}
	}

	// Test listing with prefix
	files, err := mockClient.ListFiles(context.Background(), "uploads/user1/")
	if err != nil {
		t.Errorf("ListFiles() error = %v", err)
	}

	expectedCount := 2
	if len(files) != expectedCount {
		t.Errorf("ListFiles() returned %d files, want %d", len(files), expectedCount)
	}

	// Verify correct files are returned
	fileMap := make(map[string]bool)
	for _, file := range files {
		fileMap[file] = true
	}

	if !fileMap["uploads/user1/file1.txt"] || !fileMap["uploads/user1/file2.pdf"] {
		t.Error("ListFiles() did not return expected files")
	}
}
