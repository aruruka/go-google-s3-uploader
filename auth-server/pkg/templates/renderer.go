package templates

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"time"
)

// TemplateRendererIface defines the interface for template rendering
type TemplateRendererIface interface {
	RenderTemplate(w io.Writer, name string, data interface{}) error
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

// loadTemplates loads templates from the file system (not embedded for now)
func (tr *TemplateRenderer) loadTemplates() error {
	// For now, we'll create templates directly without embed
	// This allows us to focus on the handlers first
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
func (tr *TemplateRenderer) RenderTemplate(w io.Writer, name string, data interface{}) error {
	// For initial implementation, we'll return hardcoded templates
	// This will be replaced with actual template parsing later
	switch name {
	case "login.html":
		return tr.renderLoginPage(w, data)
	case "callback.html":
		return tr.renderCallbackPage(w, data)
	case "error.html":
		return tr.renderErrorPage(w, data)
	default:
		return fmt.Errorf("template %s not found", name)
	}
}

// Temporary hardcoded templates for initial testing
func (tr *TemplateRenderer) renderLoginPage(w io.Writer, data interface{}) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - Google S3 Uploader</title>
    <link href="/static/css/styles.css" rel="stylesheet">
</head>
<body>
    <div class="container">
        <div class="auth-container">
            <h1>Welcome to Google S3 Uploader</h1>
            <p>Please sign in with your Google account to continue</p>
            <a href="/auth/google" class="google-signin-btn">
                <svg width="18" height="18" viewBox="0 0 18 18">
                    <path fill="#4285f4" d="m18 9.2c0-.7-.1-1.4-.2-2h-8.8v3.9h5.1c-.2 1.1-.9 2-1.8 2.7v2.2h2.9c1.7-1.6 2.8-3.9 2.8-6.8z"/>
                    <path fill="#34a853" d="m9 18c2.4 0 4.5-.8 6-2.2l-2.9-2.2c-.8.6-1.9.9-3.1.9-2.4 0-4.4-1.6-5.1-3.9h-3v2.3c1.6 3.1 4.7 5.1 8.1 5.1z"/>
                    <path fill="#fbbc04" d="m3.9 10.7c-.2-.6-.2-1.2 0-1.8v-2.2h-3c-.7 1.4-.7 3.1 0 4.5l3-2.5z"/>
                    <path fill="#ea4335" d="m9 3.6c1.3 0 2.5.4 3.4 1.3l2.5-2.5c-1.5-1.4-3.5-2.4-5.9-2.4-3.4 0-6.5 2-8.1 5.1l3 2.3c.7-2.3 2.7-3.8 5.1-3.8z"/>
                </svg>
                Sign in with Google
            </a>
        </div>
    </div>
</body>
</html>`
	_, err := w.Write([]byte(html))
	return err
}

func (tr *TemplateRenderer) renderCallbackPage(w io.Writer, data interface{}) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Authentication Success</title>
    <link href="/static/css/styles.css" rel="stylesheet">
</head>
<body>
    <div class="container">
        <div class="auth-container">
            <h1>Authentication Successful!</h1>
            <p>Redirecting to application...</p>
            <script>
                setTimeout(() => {
                    window.location.href = '{{.AppServerURL}}';
                }, 2000);
            </script>
        </div>
    </div>
</body>
</html>`
	_, err := w.Write([]byte(html))
	return err
}

func (tr *TemplateRenderer) renderErrorPage(w io.Writer, data interface{}) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error</title>
    <link href="/static/css/styles.css" rel="stylesheet">
</head>
<body>
    <div class="container">
        <div class="auth-container error">
            <h1>Authentication Error</h1>
            <p>Sorry, there was an error during authentication. Please try again.</p>
            <a href="/login" class="retry-btn">Try Again</a>
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
