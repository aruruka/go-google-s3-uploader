package templates

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"time"

	"app-server/pkg/models"
)

// TemplateRendererIface defines the interface for template rendering
type TemplateRendererIface interface {
	RenderTemplate(w io.Writer, name string, data any) error
}

// TemplateRenderer implements the template rendering functionality
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new template renderer instance
func NewTemplateRenderer() (TemplateRendererIface, error) {
	renderer := &TemplateRenderer{}
	if err := renderer.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}
	return renderer, nil
}

// loadTemplates loads templates (for now, using hardcoded templates)
func (tr *TemplateRenderer) loadTemplates() error {
	tr.templates = template.New("").Funcs(template.FuncMap{
		"formatDate":     formatDate,
		"formatFileSize": formatFileSize,
		"json":           toJSON,
		"safe":           safe,
		"dict":           dict,
	})

	return nil
}

// RenderTemplate renders a template with the given data
func (tr *TemplateRenderer) RenderTemplate(w io.Writer, name string, data any) error {
	// For initial implementation, we'll return hardcoded templates
	switch name {
	case "home.html":
		return tr.renderHomePage(w, data)
	case "upload.html":
		return tr.renderUploadPage(w, data)
	case "success.html":
		return tr.renderSuccessPage(w, data)
	case "error.html":
		return tr.renderErrorPage(w, data)
	default:
		return fmt.Errorf("template %s not found", name)
	}
}

// Temporary hardcoded templates for initial testing
func (tr *TemplateRenderer) renderHomePage(w io.Writer, _ any) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Google S3 Uploader</title>
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <!-- Header Component - Authenticated State -->
    <header class="header">
        <nav class="navbar">
            <div class="nav-container">
                <div class="nav-brand">
                    <h1>🚀 Google S3 Uploader</h1>
                </div>
                <div class="nav-menu">
                    <div class="nav-user">
                        <span class="user-info">👋 Hello, John Doe!</span>
                        <a href="/logout" class="nav-link">Logout</a>
                    </div>
                </div>
            </div>
        </nav>
    </header>

    <main class="main-content">
        <!-- Success Flash Message -->
        <div class="flash-message flash-success">
            ✅ Welcome! You have successfully logged in.
        </div>

        <!-- Home Page Content -->
        <div class="home-container">
            <div class="welcome-section">
                <h1 class="app-title">📱 Google S3 Uploader</h1>
                <p class="app-description">
                    Upload your images to AWS S3 using secure Google authentication.
                    A modern, cloud-native application built with Go.
                </p>
            </div>

            <!-- Authenticated User Content -->
            <div class="action-card">
                <h2>🚀 Ready to Upload</h2>
                <p>Welcome back, John Doe! You're authenticated and ready to upload files.</p>
                <a href="/upload" class="upload-btn">📷 Go to Upload Page</a>
            </div>

            <div class="action-card">
                <h2>✨ Features</h2>
                <ul class="feature-list">
                    <li>🔒 Secure Google OAuth 2.0 authentication</li>
                    <li>☁️ Direct upload to AWS S3</li>
                    <li>🖼️ Support for multiple image formats</li>
                    <li>📱 Responsive design</li>
                    <li>⚡ Fast and lightweight</li>
                </ul>
            </div>
        </div>
    </main>

    <!-- Footer Component -->
    <footer class="footer">
        <div class="footer-container">
            <div class="footer-content">
                <p>&copy; 2025 Google S3 Uploader. Built with Go 💙</p>
                <div class="footer-links">
                    <a href="https://golang.org" target="_blank">Go Lang</a>
                    <a href="https://aws.amazon.com/s3/" target="_blank">AWS S3</a>
                    <a href="https://developers.google.com/identity" target="_blank">Google OAuth</a>
                </div>
            </div>
        </div>
    </footer>

    <script src="/static/js/app.js"></script>
</body>
</html>`
	_, err := w.Write([]byte(html))
	return err
}

func (tr *TemplateRenderer) renderUploadPage(w io.Writer, _ any) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Upload File - Google S3 Uploader</title>
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <!-- Header Component - Authenticated State -->
    <header class="header">
        <nav class="navbar">
            <div class="nav-container">
                <div class="nav-brand">
                    <h1>🚀 Google S3 Uploader</h1>
                </div>
                <div class="nav-menu">
                    <div class="nav-user">
                        <span class="user-info">👋 Hello, John Doe!</span>
                        <a href="/logout" class="nav-link">Logout</a>
                    </div>
                </div>
            </div>
        </nav>
    </header>

    <main class="main-content">
        <!-- Upload Page Content -->
        <div class="upload-container">
            <div class="upload-card">
                <h1 class="upload-title">📷 Upload Image to S3</h1>
                
                <div class="upload-info">
                    <strong>ℹ️ Upload Information:</strong>
                    <ul style="margin: 0.5rem 0 0 1rem;">
                        <li>Supported formats: JPG, PNG, GIF, WebP</li>
                        <li>Maximum file size: 10 MB</li>
                        <li>Files will be stored securely in AWS S3</li>
                    </ul>
                </div>

                <form action="/api/upload" method="post" enctype="multipart/form-data" class="upload-form" id="uploadForm">
                    <div class="form-group">
                        <label for="file" class="form-label">Choose image file:</label>
                        <input type="file" id="file" name="file" accept="image/*" required class="file-input">
                        <div class="file-preview" id="filePreview">
                            <img id="previewImage" class="preview-image" alt="Preview">
                            <p id="fileName"></p>
                        </div>
                    </div>
                    
                    <div class="progress-bar" id="progressBar">
                        <div class="progress-fill" id="progressFill"></div>
                    </div>
                    
                    <button type="submit" class="upload-btn" id="uploadBtn">
                        🚀 Upload to S3
                    </button>
                </form>
            </div>
        </div>
    </main>

    <!-- Footer Component -->
    <footer class="footer">
        <div class="footer-container">
            <div class="footer-content">
                <p>&copy; 2025 Google S3 Uploader. Built with Go 💙</p>
                <div class="footer-links">
                    <a href="https://golang.org" target="_blank">Go Lang</a>
                    <a href="https://aws.amazon.com/s3/" target="_blank">AWS S3</a>
                    <a href="https://developers.google.com/identity" target="_blank">Google OAuth</a>
                </div>
            </div>
        </div>
    </footer>

    <script src="/static/js/app.js"></script>
    <script src="/static/js/upload.js"></script>
</body>
</html>`
	_, err := w.Write([]byte(html))
	return err
}

func (tr *TemplateRenderer) renderSuccessPage(w io.Writer, data interface{}) error {
	// Cast data to PageData
	pageData, ok := data.(*models.PageData)
	if !ok {
		return fmt.Errorf("invalid data type for success page")
	}

	// Cast page data to SuccessData
	successData, ok := pageData.Data.(*models.SuccessData)
	if !ok {
		return fmt.Errorf("invalid success data type")
	}

	// Get user info
	userName := "User"
	if pageData.User != nil {
		userName = pageData.User.Name
	}

	// Get file details
	upload := successData.Upload
	if upload == nil {
		return fmt.Errorf("no upload data provided")
	}

	// Format file size
	fileSize := formatFileSize(upload.Size)

	// Format upload time
	uploadTime := formatDate(upload.UploadedAt)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <link href="/static/css/style.css" rel="stylesheet">
    <style>
        /* Inline styles from success.html template */
        .success-container {
            max-width: 600px;
            margin: 2rem auto;
            padding: 2rem;
            text-align: center;
        }
        .success-card {
            background: white;
            border-radius: 8px;
            padding: 2rem;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .success-icon {
            font-size: 4rem;
            margin-bottom: 1rem;
        }
        .success-title {
            color: #28a745;
            margin-bottom: 1rem;
        }
        .file-info {
            background: #f8f9fa;
            border-radius: 6px;
            padding: 1rem;
            margin: 1rem 0;
        }
        .file-url {
            word-break: break-all;
            background: #e9ecef;
            padding: 0.5rem;
            border-radius: 4px;
            font-family: monospace;
            margin: 0.5rem 0;
        }
        .action-buttons {
            display: flex;
            gap: 1rem;
            justify-content: center;
            margin-top: 2rem;
        }
        .btn {
            padding: 10px 20px;
            border-radius: 6px;
            text-decoration: none;
            font-weight: 500;
            transition: background 0.2s;
        }
        .btn-primary {
            background: #007bff;
            color: white;
        }
        .btn-primary:hover {
            background: #0056b3;
        }
        .btn-secondary {
            background: #6c757d;
            color: white;
        }
        .btn-secondary:hover {
            background: #545b62;
        }
    </style>
</head>
<body>
    <!-- Header Component - Authenticated State -->
    <header class="header">
        <nav class="navbar">
            <div class="nav-container">
                <div class="nav-brand">
                    <h1>� Google S3 Uploader</h1>
                </div>
                <div class="nav-menu">
                    <div class="nav-user">
                        <span class="user-info">👋 Hello, %s!</span>
                        <a href="/logout" class="nav-link">Logout</a>
                    </div>
                </div>
            </div>
        </nav>
    </header>

    <main class="main-content">
        <!-- Success Flash Message -->
        <div class="flash-message flash-success">
            🎉 File uploaded successfully to AWS S3!
        </div>

        <!-- Success Page Content -->
        <div class="success-container">
            <div class="success-card">
                <div class="success-icon">🎉</div>
                <h1 class="success-title">Upload Successful!</h1>
                
                <div class="file-info">
                    <h3>📁 File Details</h3>
                    <p><strong>Filename:</strong> %s</p>
                    <p><strong>Size:</strong> %s</p>
                    <p><strong>Type:</strong> %s</p>
                    <p><strong>Uploaded:</strong> %s</p>
                </div>
                
                <div class="file-info">
                    <h3>🔗 File URL</h3>
                    <div class="file-url">%s</div>
                    <button onclick="copyToClipboard('%s')" class="btn btn-secondary">
                        📋 Copy URL
                    </button>
                </div>
                
                <div class="action-buttons">
                    <a href="/upload" class="btn btn-primary">� Upload Another</a>
                    <a href="/" class="btn btn-secondary">🏠 Go Home</a>
                </div>
            </div>
        </div>
    </main>

    <!-- Footer Component -->
    <footer class="footer">
        <div class="footer-container">
            <div class="footer-content">
                <p>&copy; 2025 Google S3 Uploader. Built with Go 💙</p>
                <div class="footer-links">
                    <a href="https://golang.org" target="_blank">Go Lang</a>
                    <a href="https://aws.amazon.com/s3/" target="_blank">AWS S3</a>
                    <a href="https://developers.google.com/identity" target="_blank">Google OAuth</a>
                </div>
            </div>
        </div>
    </footer>

    <script>
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(function() {
        alert('✅ URL copied to clipboard!');
    }, function(err) {
        console.error('Could not copy text: ', err);
        // Fallback for older browsers
        const textArea = document.createElement("textarea");
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        try {
            document.execCommand('copy');
            alert('✅ URL copied to clipboard!');
        } catch (err) {
            alert('❌ Failed to copy URL');
        }
        document.body.removeChild(textArea);
    });
}
    </script>
</body>
</html>`, pageData.Title, userName, upload.Filename, fileSize, upload.ContentType, uploadTime, upload.S3URL, upload.S3URL)

	_, err := w.Write([]byte(html))
	return err
}

func (tr *TemplateRenderer) renderErrorPage(w io.Writer, data interface{}) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error - Google S3 Uploader</title>
    <link href="/static/css/style.css" rel="stylesheet">
</head>
<body>
    <div class="container">
        <div class="error-container">
            <div class="error-icon">❌</div>
            <h1>Upload Error</h1>
            <p>Sorry, there was an error uploading your file. Please try again.</p>
            
            <div class="error-actions">
                <a href="/upload" class="btn primary">Try Again</a>
                <a href="/" class="btn secondary">Back to Home</a>
            </div>
        </div>
    </div>
</body>
</html>`
	_, err := w.Write([]byte(html))
	return err
}

// Helper functions for templates

// formatDate formats a time.Time to a readable string
func formatDate(t time.Time) string {
	return t.Format("January 2, 2006 at 3:04 PM")
}

// formatFileSize formats bytes to human readable format
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// toJSON converts data to JSON string
func toJSON(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(b)
}

// safe marks a string as safe for HTML output
func safe(s string) template.HTML {
	return template.HTML(s)
}

// dict creates a map for use in templates
func dict(values ...interface{}) map[string]interface{} {
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		if i+1 < len(values) {
			dict[fmt.Sprintf("%v", values[i])] = values[i+1]
		}
	}
	return dict
}
