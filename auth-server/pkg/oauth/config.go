package oauth

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

func NewConfig() (*Config, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("REDIRECT_URL")
	isProduction := os.Getenv("ENV") == "production"

	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID environment variable is required")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET environment variable is required")
	}
	if redirectURL == "" {
		if isProduction {
			return nil, fmt.Errorf("REDIRECT_URL environment variable is required in production")
		}
		redirectURL = "http://localhost:8081/auth/callback"
		log.Printf("Using default redirect URL: %s", redirectURL)
	}

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to get OIDC provider: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	return &Config{
		OAuth2Config: oauth2Config,
		Verifier:     verifier,
	}, nil
}

func (c *Config) GetAuthURL(state string) string {
	return c.OAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (c *Config) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.OAuth2Config.Exchange(ctx, code)
}

func (c *Config) VerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	return c.Verifier.Verify(ctx, rawIDToken)
}
