package oteltrace

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{} `json:"-"`

		CollectorEndpoint string `json:"collector_endpoint,omitempty" toml:"collector_endpoint,omitempty" env:"OTELGRPC_COLLECTOR_ENDPOINT"`
		Insecure          bool   `json:"insecure,omitempty"           toml:"insecure,omitempty"           env:"OTELGRPC_INSECURE"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.CollectorEndpoint, validation.Required),
	)
}
