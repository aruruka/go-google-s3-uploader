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

// ErrorData represents data for error pages
type ErrorData struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}
