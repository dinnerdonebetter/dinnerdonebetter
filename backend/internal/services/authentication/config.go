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

		ClientID     string `env:"CLIENT_ID"     json:"clientID,omitempty"`
		ClientSecret string `env:"CLIENT_SECRET" json:"clientSecret,omitempty"`
		CallbackURL  string `env:"CALLBACK_URL"  json:"callbackURL,omitempty"`
	}

	SSOConfigs struct {
		Google GoogleSSOConfig `envPrefix:"GOOGLE_" json:"google,omitempty"`
	}

	// Config represents our passwords configuration.
	Config struct {
		_ struct{} `json:"-"`

		SSO                   SSOConfigs    `envPrefix:"SSO_CONFIG_"       json:"sso,omitempty"`
		DataChangesTopicName  string        `env:"DATA_CHANGES_TOPIC_NAME" json:"dataChanges,omitempty"`
		JWTAudience           string        `env:"JWT_AUDIENCE"            json:"jwtAudience,omitempty"`
		JWTSigningKey         string        `env:"JWT_SIGNING_KEY"         json:"jwtSigningKey"`
		OAuth2                OAuth2Config  `envPrefix:"OAUTH2"            json:"oauth2,omitempty"`
		JWTLifetime           time.Duration `env:"JWT_LIFETIME"            json:"jwtLifetime"`
		Debug                 bool          `env:"DEBUG"                   json:"debug,omitempty" envDefault:"false"`
		EnableUserSignup      bool          `env:"ENABLE_USER_SIGNUP"      json:"enableUserSignup,omitempty"`
		MinimumUsernameLength uint8         `env:"MINIMUM_USERNAME_LENGTH" json:"minimumUsernameLength,omitempty" envDefault:"4"`
		MinimumPasswordLength uint8         `env:"MINIMUM_PASSWORD_LENGTH" json:"minimumPasswordLength,omitempty" envDefault:"8"`
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

	Domain               string        `env:"DOMAIN"                 json:"domain"`
	AccessTokenLifespan  time.Duration `env:"ACCESS_TOKEN_LIFESPAN"  json:"accessTokenLifespan"`
	RefreshTokenLifespan time.Duration `env:"REFRESH_TOKEN_LIFESPAN" json:"refreshTokenLifespan"`
	Debug                bool          `env:"DEBUG"                  json:"debug" envDefault:"false"`
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
