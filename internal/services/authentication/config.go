package authentication

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// DefaultCookieName is the default Cookie.BucketName.
	DefaultCookieName = "ddb_api_cookie"
	// DefaultCookieLifetime is the how long a cookie is valid.
	DefaultCookieLifetime = 24 * time.Hour

	staticError             = "error encountered, please try again later"
	pasetoKeyRequiredLength = 32
	pasetoDataKey           = "paseto_data"
	maxPASETOLifetime       = 10 * time.Minute
)

type (
	// CookieConfig holds our cookie settings.
	CookieConfig struct {
		_ struct{}

		Name       string        `json:"name,omitempty"       mapstructure:"name"        toml:"name,omitempty"`
		Domain     string        `json:"domain,omitempty"     mapstructure:"domain"      toml:"domain,omitempty"`
		HashKey    string        `json:"hashKey,omitempty"    mapstructure:"hash_key"    toml:"hash_key,omitempty"`
		BlockKey   string        `json:"blockKey,omitempty"   mapstructure:"signing_key" toml:"signing_key,omitempty"`
		Lifetime   time.Duration `json:"lifetime,omitempty"   mapstructure:"lifetime"    toml:"lifetime,omitempty"`
		SecureOnly bool          `json:"secureOnly,omitempty" mapstructure:"secure_only" toml:"secure_only,omitempty"`
	}

	// PASETOConfig holds our PASETO settings.
	PASETOConfig struct {
		_ struct{}

		Issuer       string        `json:"issuer,omitempty"       mapstructure:"issuer"         toml:"issuer,omitempty"`
		LocalModeKey []byte        `json:"localModeKey,omitempty" mapstructure:"local_mode_key" toml:"local_mode_key,omitempty"`
		Lifetime     time.Duration `json:"lifetime,omitempty"     mapstructure:"lifetime"       toml:"lifetime,omitempty"`
	}

	// Config represents our passwords configuration.
	Config struct {
		_ struct{}

		DataChangesTopicName  string       `json:"dataChanges,omitempty"           mapstructure:"data_changes"            toml:"data_changes,omitempty"`
		Cookies               CookieConfig `json:"cookies,omitempty"               mapstructure:"cookies"                 toml:"cookies,omitempty"`
		PASETO                PASETOConfig `json:"paseto,omitempty"                mapstructure:"paseto"                  toml:"paseto,omitempty"`
		Debug                 bool         `json:"debug,omitempty"                 mapstructure:"debug"                   toml:"debug,omitempty"`
		EnableUserSignup      bool         `json:"enableUserSignup,omitempty"      mapstructure:"enable_user_signup"      toml:"enable_user_signup,omitempty"`
		MinimumUsernameLength uint8        `json:"minimumUsernameLength,omitempty" mapstructure:"minimum_username_length" toml:"minimum_username_length,omitempty"`
		MinimumPasswordLength uint8        `json:"minimumPasswordLength,omitempty" mapstructure:"minimum_password_length" toml:"minimum_password_length,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*CookieConfig)(nil)

// ValidateWithContext validates a CookieConfig struct.
func (cfg *CookieConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Name, validation.Required),
		validation.Field(&cfg.Domain, validation.Required),
		validation.Field(&cfg.Lifetime, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*PASETOConfig)(nil)

// ValidateWithContext validates a PASETOConfig struct.
func (cfg *PASETOConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Issuer, validation.Required),
		validation.Field(&cfg.LocalModeKey, validation.Required, validation.Length(pasetoKeyRequiredLength, pasetoKeyRequiredLength)),
	)
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Cookies, validation.Required),
		validation.Field(&cfg.PASETO, validation.Required),
		validation.Field(&cfg.MinimumUsernameLength, validation.Required),
		validation.Field(&cfg.MinimumPasswordLength, validation.Required),
	)
}
