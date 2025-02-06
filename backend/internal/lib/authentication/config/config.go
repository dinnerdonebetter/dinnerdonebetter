package authcfg

import (
	"time"

	tokenscfg "github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens/config"
)

type (
	GoogleSSOConfig struct {
		_ struct{} `json:"-"`

		ClientID     string `env:"CLIENT_ID"     json:"clientID,omitempty"`
		ClientSecret string `env:"CLIENT_SECRET" json:"clientSecret,omitempty"`
		CallbackURL  string `env:"CALLBACK_URL"  json:"callbackURL,omitempty"`
	}

	SSOConfigs struct {
		Google GoogleSSOConfig `envPrefix:"GOOGLE_SSO_" json:"google,omitempty"`
	}

	// Config is our configuration.
	Config struct {
		_ struct{} `json:"-"`

		Tokens                  tokenscfg.Config `envPrefix:"TOKENS_"              json:"tokens"`
		SSO                     SSOConfigs       `envPrefix:"SSO_CONFIG_"          json:"sso,omitempty"`
		MaxAccessTokenLifetime  time.Duration    `env:"MAX_ACCESS_TOKEN_LIFETIME"  json:"maxAccessTokenLifetime"`
		MaxRefreshTokenLifetime time.Duration    `env:"MAX_REFRESH_TOKEN_LIFETIME" json:"maxRefreshTokenLifetime"`
		Debug                   bool             `env:"DEBUG"                      json:"debug,omitempty"`
		EnableUserSignup        bool             `env:"ENABLE_USER_SIGNUP"         json:"enableUserSignup,omitempty"`
		MinimumUsernameLength   uint8            `env:"MINIMUM_USERNAME_LENGTH"    json:"minimumUsernameLength,omitempty"`
		MinimumPasswordLength   uint8            `env:"MINIMUM_PASSWORD_LENGTH"    json:"minimumPasswordLength,omitempty"`
	}
)
