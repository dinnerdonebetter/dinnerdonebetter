package authentication

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	staticError = "error encountered, please try again later"
)

type (
	GoogleSSOConfig struct {
		_ struct{} `json:"-"`

		ClientID     string `json:"clientID,omitempty"     toml:"client_id,omitempty"`
		ClientSecret string `json:"clientSecret,omitempty" toml:"client_secret,omitempty"`
		CallbackURL  string `json:"callbackURL,omitempty"  toml:"callback_url,omitempty"`
	}

	SSOConfigs struct {
		Google GoogleSSOConfig `json:"google,omitempty" toml:"google,omitempty"`
	}

	// Config represents our passwords configuration.
	Config struct {
		_ struct{} `json:"-"`

		SSO                   SSOConfigs    `json:"sso,omitempty"                   toml:"sso,omitempty"`
		DataChangesTopicName  string        `json:"dataChanges,omitempty"           toml:"data_changes,omitempty"`
		JWTAudience           string        `json:"jwtAudience,omitempty"           toml:"jwt_audience,omitempty"`
		JWTSigningKey         string        `json:"jwtSigningKey"                   toml:"jwt_signing_key"`
		OAuth2                OAuth2Config  `json:"oauth2,omitempty"                toml:"oauth2,omitempty"`
		JWTLifetime           time.Duration `json:"jwtLifetime"                     toml:"jwt_lifetime"`
		Debug                 bool          `json:"debug,omitempty"                 toml:"debug,omitempty"`
		EnableUserSignup      bool          `json:"enableUserSignup,omitempty"      toml:"enable_user_signup,omitempty"`
		MinimumUsernameLength uint8         `json:"minimumUsernameLength,omitempty" toml:"minimum_username_length,omitempty"`
		MinimumPasswordLength uint8         `json:"minimumPasswordLength,omitempty" toml:"minimum_password_length,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.MinimumUsernameLength, validation.Required),
		validation.Field(&cfg.MinimumPasswordLength, validation.Required),
	)
}

// OAuth2Config represents our database configuration.
type OAuth2Config struct {
	_ struct{} `json:"-"`

	Domain               string        `json:"domain"               toml:"domain,omitempty"`
	AccessTokenLifespan  time.Duration `json:"accessTokenLifespan"  toml:"access_token_lifespan,omitempty"`
	RefreshTokenLifespan time.Duration `json:"refreshTokenLifespan" toml:"refresh_token_lifespan,omitempty"`
	Debug                bool          `json:"debug"                toml:"debug,omitempty"`
}

var _ validation.ValidatableWithContext = (*OAuth2Config)(nil)

// ValidateWithContext validates a OAuth2Config struct.
func (cfg OAuth2Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.AccessTokenLifespan, validation.Required),
		validation.Field(&cfg.RefreshTokenLifespan, validation.Required),
		validation.Field(&cfg.Domain, validation.Required),
	)
}
