package redis

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures a Redis-backed consumer.
type Config struct {
	Username       string   `json:"username" mapstructure:"username" toml:"username,omitempty"`
	Password       string   `json:"password,omitempty" mapstructure:"password" toml:"password,omitempty"`
	QueueAddresses []string `json:"queueAddress" mapstructure:"queue_address" toml:"queue_address,omitempty"`
	DB             int      `json:"database,omitempty" mapstructure:"database" toml:"database,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.QueueAddresses, validation.Required),
	)
}
