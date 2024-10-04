package config

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config represents our database configuration.
	Config struct {
		_ struct{} `json:"-"`

		OAuth2TokenEncryptionKey string        `json:"oauth2TokenEncryptionKey" toml:"oauth2_token_encryption_key,omitempty"`
		ConnectionDetails        string        `json:"connectionDetails"        toml:"connection_details,omitempty"`
		Debug                    bool          `json:"debug"                    toml:"debug,omitempty"`
		LogQueries               bool          `json:"logQueries"               toml:"log_queries,omitempty"`
		RunMigrations            bool          `json:"runMigrations"            toml:"run_migrations,omitempty"`
		MaxPingAttempts          uint64        `json:"maxPingAttempts"          toml:"max_ping_attempts,omitempty"`
		PingWaitPeriod           time.Duration `json:"pingWaitPeriod"           toml:"ping_wait_period,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates an DatabaseSettings struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.ConnectionDetails, validation.Required),
		validation.Field(&cfg.OAuth2TokenEncryptionKey, validation.Required),
	)
}
