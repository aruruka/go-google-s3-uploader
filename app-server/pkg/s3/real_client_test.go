package s3

import (
	"os"
	"strings"
	"testing"
)

// TestNewS3Client tests the S3 client constructor
func TestNewS3Client_Constructor(t *testing.T) {
	tests := []struct {
		name        string
		bucketName  string
		region      string
		expectError bool
	}{
		{
			name:        "valid config",
			bucketName:  "test-bucket",
			region:      "us-west-2",
			expectError: false,
		},
		{
			name:        "missing bucket name",
			bucketName:  "",
			region:      "us-west-2",
			expectError: true,
		},
		{
			name:        "default region when empty",
			bucketName:  "test-bucket",
			region:      "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.bucketName != "" {
				os.Setenv("S3_BUCKET_NAME", tt.bucketName)
			} else {
				os.Unsetenv("S3_BUCKET_NAME")
			}

			if tt.region != "" {
				os.Setenv("AWS_REGION", tt.region)
			} else {
				os.Unsetenv("AWS_REGION")
			}

			// Clean up after test
			defer func() {
				os.Unsetenv("S3_BUCKET_NAME")
				os.Unsetenv("AWS_REGION")
			}()

			_, err := NewS3Client()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				// Skip AWS config errors in unit tests
				if !strings.Contains(err.Error(), "failed to load AWS config") {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestS3Client_GetFileURL_Implementation(t *testing.T) {
	client := &S3Client{
		bucketName: "my-test-bucket",
		region:     "us-west-2",
	}

	tests := []struct {
		name        string
		key         string
		expectedURL string
	}{
		{
			name:        "simple file",
			key:         "file.txt",
			expectedURL: "https://my-test-bucket.s3.us-west-2.amazonaws.com/file.txt",
		},
		{
			name:        "nested path",
			key:         "uploads/user123/image.jpg",
			expectedURL: "https://my-test-bucket.s3.us-west-2.amazonaws.com/uploads/user123/image.jpg",
		},
		{
			name:        "file with spaces",
			key:         "documents/my file.pdf",
			expectedURL: "https://my-test-bucket.s3.us-west-2.amazonaws.com/documents/my file.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualURL := client.GetFileURL(tt.key)
			if actualURL != tt.expectedURL {
				t.Errorf("Expected URL %s, got %s", tt.expectedURL, actualURL)
			}
		})
	}
}

// TestS3Client_ErrorHandling tests error scenarios
func TestS3Client_ErrorHandling(t *testing.T) {
	// Test environment variable validation
	t.Run("missing bucket name", func(t *testing.T) {
		os.Unsetenv("S3_BUCKET_NAME")
		defer os.Setenv("S3_BUCKET_NAME", "test-bucket")

		_, err := NewS3Client()
		if err == nil {
			t.Error("Expected error for missing bucket name")
		}
		if !strings.Contains(err.Error(), "S3_BUCKET_NAME environment variable is required") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}

// Benchmark tests for performance
func BenchmarkS3Client_GetFileURL_Real(b *testing.B) {
	client := &S3Client{
		bucketName: "benchmark-bucket",
		region:     "us-east-1",
	}

	key := "test/benchmark/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.GetFileURL(key)
	}
}

// Integration test helper functions
func TestS3Client_IntegrationHelper(t *testing.T) {
	// Skip integration tests in unit test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires real AWS credentials
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		t.Skip("S3_BUCKET_NAME not set, skipping integration test")
	}

	client, err := NewS3Client()
	if err != nil {
		t.Skipf("Cannot create S3 client: %v", err)
	}

	// Test file URL generation
	testKey := "integration-test/unit-test-file.txt"
	url := client.GetFileURL(testKey)

	expectedPrefix := "https://" + bucketName + ".s3."
	if !strings.HasPrefix(url, expectedPrefix) {
		t.Errorf("URL should start with %s, got %s", expectedPrefix, url)
	}
}

// Test data validation
func TestS3Client_DataValidation(t *testing.T) {
	client := &S3Client{
		bucketName: "test-bucket",
		region:     "us-east-1",
	}

	tests := []struct {
		name string
		key  string
		want bool // true if key should be valid
	}{
		{"empty key", "", false},
		{"normal key", "file.txt", true},
		{"nested key", "folder/subfolder/file.txt", true},
		{"key with unicode", "文件.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := client.GetFileURL(tt.key)
			hasKey := strings.Contains(url, tt.key)

			if tt.want && !hasKey {
				t.Errorf("Expected URL to contain key %s, got %s", tt.key, url)
			}
		})
	}
}
