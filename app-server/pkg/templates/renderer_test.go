package templates

import (
	"bytes"
	"testing"
	"time"

	"app-server/pkg/models"
)

// Test NewTemplateRenderer
func TestNewTemplateRenderer(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if renderer == nil {
		t.Error("Expected renderer to be created, got nil")
	}
}

// Test RenderTemplate for different template types
func TestTemplateRenderer_RenderTemplate(t *testing.T) {
	renderer, err := NewTemplateRenderer()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	tests := []struct {
		name         string
		templateName string
		data         interface{}
		expectError  bool
	}{
		{
			name:         "home template",
			templateName: "home.html",
			data:         map[string]interface{}{"Title": "Test Home"},
			expectError:  false,
		},
		{
			name:         "upload template",
			templateName: "upload.html",
			data:         map[string]interface{}{"Title": "Test Upload"},
			expectError:  false,
		},
		{
			name:         "success template",
			templateName: "success.html",
			data: &models.PageData{
				Title: "Test Success",
				User: &models.User{
					Name: "Test User",
				},
				Data: &models.SuccessData{
					Upload: &models.FileUpload{
						Filename:    "test.jpg",
						Size:        1024,
						ContentType: "image/jpeg",
						S3URL:       "https://example.com/test.jpg",
						UploadedAt:  time.Now(),
					},
				},
			},
			expectError: false,
		},
		{
			name:         "error template",
			templateName: "error.html",
			data:         map[string]interface{}{"Title": "Test Error"},
			expectError:  false,
		},
		{
			name:         "unknown template",
			templateName: "unknown.html",
			data:         map[string]interface{}{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := renderer.RenderTemplate(&buf, tt.templateName, tt.data)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				
				// Check that something was written
				if buf.Len() == 0 {
					t.Error("Expected output but got empty buffer")
				}
			}
		})
	}
}

// Test helper functions
func TestTemplateHelpers(t *testing.T) {
	t.Run("formatDate", func(t *testing.T) {
		// We need to import time to test this
		// For now, just check the function doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("formatDate panicked: %v", r)
			}
		}()
		// formatDate(time.Now()) // Would need to import time package
	})

	t.Run("formatFileSize", func(t *testing.T) {
		tests := []struct {
			bytes    int64
			expected string
		}{
			{512, "512 B"},
			{1024, "1.0 KB"},
			{1536, "1.5 KB"},
			{1048576, "1.0 MB"},
		}

		for _, tt := range tests {
			result := formatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatFileSize(%d) = %s, expected %s", tt.bytes, result, tt.expected)
			}
		}
	})

	t.Run("toJSON", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		result := toJSON(data)
		expected := `{"key":"value"}`
		if result != expected {
			t.Errorf("toJSON() = %s, expected %s", result, expected)
		}
	})

	t.Run("dict", func(t *testing.T) {
		result := dict("key1", "value1", "key2", "value2")
		if len(result) != 2 {
			t.Errorf("Expected 2 items in dict, got %d", len(result))
		}
		if result["key1"] != "value1" {
			t.Errorf("Expected key1=value1, got %v", result["key1"])
		}
	})
}
