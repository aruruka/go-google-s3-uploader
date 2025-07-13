package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"auth-server/pkg/models"
	"auth-server/pkg/oauth"
	"auth-server/pkg/templates"
)

// AuthHandlerIface defines the interface for authentication handlers
type AuthHandlerIface interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleGoogleAuth(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
}

// AuthHandler implements authentication handlers
type AuthHandler struct {
	oauthConfig *oauth.Config
	renderer    templates.TemplateRendererIface
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(oauthConfig *oauth.Config, renderer templates.TemplateRendererIface) AuthHandlerIface {
	return &AuthHandler{
		oauthConfig: oauthConfig,
		renderer:    renderer,
	}
}

// HandleLogin displays the login page
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	pageData := &models.PageData{
		Title: "Login - Google S3 Uploader",
		Data: &models.LoginData{
			RedirectURL: r.URL.Query().Get("redirect"),
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.RenderTemplate(w, "login.html", pageData); err != nil {
		log.Printf("Failed to render login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleGoogleAuth redirects to Google OAuth
func (h *AuthHandler) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	// Generate state token
	state, err := generateStateToken()
	if err != nil {
		log.Printf("Failed to generate state token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store state in session (for now, we'll use a cookie)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Get auth URL and redirect
	authURL := h.oauthConfig.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleCallback handles the OAuth callback
func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 DEBUG: HandleCallback called with URL: %s", r.URL.String())
	log.Printf("🔍 DEBUG: Query parameters: %v", r.URL.Query())

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("State cookie not found: %v", err)
		h.renderError(w, "Invalid authentication state", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		log.Printf("State mismatch: expected %s, got %s", stateCookie.Value, state)
		h.renderError(w, "Invalid authentication state", http.StatusBadRequest)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	// Handle authentication error
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		log.Printf("OAuth error: %s", errMsg)
		h.renderError(w, fmt.Sprintf("Authentication error: %s", errMsg), http.StatusBadRequest)
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Authorization code not found")
		h.renderError(w, "Authorization code not received", http.StatusBadRequest)
		return
	}

	// Exchange code for tokens
	ctx := context.Background()
	token, err := h.oauthConfig.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		h.renderError(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("No id_token field in oauth2 token")
		h.renderError(w, "Invalid token response", http.StatusInternalServerError)
		return
	}

	// Verify ID token
	idToken, err := h.oauthConfig.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		log.Printf("Failed to verify ID token: %v", err)
		h.renderError(w, "Failed to verify token", http.StatusInternalServerError)
		return
	}

	// Extract claims
	var claims struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Failed to parse claims: %v", err)
		h.renderError(w, "Failed to parse user information", http.StatusInternalServerError)
		return
	}

	// Create user model
	user := &models.User{
		ID:       idToken.Subject,
		Name:     claims.Name,
		Email:    claims.Email,
		Picture:  claims.Picture,
		Provider: "google",
		Created:  time.Now(),
	}

	// For now, just log the user information
	log.Printf("User authenticated: %s (%s)", user.Name, user.Email)

	// Store user in session (simplified - in production, use proper session management)
	userJSON, _ := json.Marshal(user)
	cookie := &http.Cookie{
		Name:     "user_session",
		Value:    base64.StdEncoding.EncodeToString(userJSON),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,                // Allow JavaScript access for debugging
		Secure:   false,                // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // Lax mode for localhost development
		Path:     "/",                  // Available on all paths
		// No domain specified - let browser handle localhost cross-port
	}

	log.Printf("🍪 Setting cookie: %s", cookie.String())
	http.SetCookie(w, cookie)

	// Instead of rendering a page, directly redirect to app-server
	log.Printf("✅ User authenticated successfully: %s (%s)", user.Name, user.Email)
	log.Printf("🔄 Redirecting to app-server...")

	// Get app server URL from environment variable, fallback to localhost:8080
	appServerURL := os.Getenv("APP_SERVER_URL")
	if appServerURL == "" {
		appServerURL = "http://localhost:8080"
	}

	http.Redirect(w, r, appServerURL, http.StatusTemporaryRedirect)
}

// HandleLogout handles user logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

// renderError renders an error page
func (h *AuthHandler) renderError(w http.ResponseWriter, message string, statusCode int) {
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

// generateStateToken generates a random state token for OAuth
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
