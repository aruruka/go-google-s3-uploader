package oauth

import (
	"context"
	"fmt"
	"log"

	"github.com/aruruka/go-google-s3-uploader/auth-server/pkg/config" // Import the new config package

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

// NewConfig initializes OAuth configuration using the provided AppConfig.
func NewConfig(appConfig *config.AppConfig) (*Config, error) {
	clientID := appConfig.GoogleClientID
	clientSecret := appConfig.GoogleClientSecret
	redirectURL := appConfig.RedirectURL

	// The validation for these variables is now handled in config.LoadConfig()
	// We just need to ensure they are not empty here, which should already be guaranteed by LoadConfig.
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is empty in AppConfig")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is empty in AppConfig")
	}
	if redirectURL == "" {
		return nil, fmt.Errorf("REDIRECT_URL is empty in AppConfig")
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

	log.Printf("OAuth Config Initialized: RedirectURL=%s", redirectURL)

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
