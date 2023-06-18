package oauth2

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config represents our database configuration.
type Config struct {
	_ struct{}

	DataChangesTopicName string        `json:"dataChangesTopicName,omitempty" toml:"data_changes_topic_name,omitempty"`
	AccessTokenLifespan  time.Duration `json:"accessTokenLifespan"            toml:"access_token_lifespan,omitempty"`
	RefreshTokenLifespan time.Duration `json:"refreshTokenLifespan"           toml:"refresh_token_lifespan,omitempty"`
	Domain               string        `json:"domain"                         toml:"domain,omitempty"`
	Debug                bool          `json:"debug"                          toml:"debug,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &cfg,
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
		validation.Field(&cfg.AccessTokenLifespan, validation.Required),
		validation.Field(&cfg.RefreshTokenLifespan, validation.Required),
		validation.Field(&cfg.Domain, validation.Required),
	)
}
