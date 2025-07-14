package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"auth-server/pkg/config" // Import the new config package
	"auth-server/pkg/models"
	"auth-server/pkg/oauth"
	"auth-server/pkg/templates"
)

type AuthHandlerIface interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleGoogleAuth(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	appConfig   *config.AppConfig // Add appConfig
	oauthConfig *oauth.Config
	renderer    templates.TemplateRendererIface
}

func NewAuthHandler(appConfig *config.AppConfig, oauthConfig *oauth.Config, renderer templates.TemplateRendererIface) AuthHandlerIface {
	return &AuthHandler{
		appConfig:   appConfig, // Store appConfig
		oauthConfig: oauthConfig,
		renderer:    renderer,
	}
}

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

func (h *AuthHandler) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := generateStateToken()
	if err != nil {
		log.Printf("Failed to generate state token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   true, // Change to true for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	authURL := h.oauthConfig.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 DEBUG: HandleCallback called with URL: %s", r.URL.String())

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

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		log.Printf("OAuth error: %s", errMsg)
		h.renderError(w, fmt.Sprintf("Authentication error: %s", errMsg), http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Authorization code not found")
		h.renderError(w, "Authorization code not received", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := h.oauthConfig.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		h.renderError(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("No id_token field in oauth2 token")
		h.renderError(w, "Invalid token response", http.StatusInternalServerError)
		return
	}

	idToken, err := h.oauthConfig.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		log.Printf("Failed to verify ID token: %v", err)
		h.renderError(w, "Failed to verify token", http.StatusInternalServerError)
		return
	}

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

	user := &models.User{
		ID:       idToken.Subject,
		Name:     claims.Name,
		Email:    claims.Email,
		Picture:  claims.Picture,
		Provider: "google",
		Created:  time.Now(),
	}

	log.Printf("User authenticated: %s (%s)", user.Name, user.Email)

	userJSON, _ := json.Marshal(user)
	cookie := &http.Cookie{
		Name:     "user_session",
		Value:    base64.StdEncoding.EncodeToString(userJSON),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   true, // Change to true for HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	log.Printf("🍪 Setting cookie: %s", cookie.String())
	http.SetCookie(w, cookie)

	log.Printf("✅ User authenticated successfully: %s (%s)", user.Name, user.Email)
	log.Printf("🔄 Redirecting to app-server...")

	// Use appConfig.AppServerURL for redirect
	http.Redirect(w, r, h.appConfig.AppServerURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

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
		http.Error(w, message, statusCode)
	}
}

func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
