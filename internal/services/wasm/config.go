package wasm

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the Service.
type Config struct {
	_ struct{}
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
	)
}
