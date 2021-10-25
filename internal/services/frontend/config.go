package frontend

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
)

// Config configures the frontend service.
type Config struct {
	_ struct{}

	Logging logging.Config `json:"logging" mapstructure:"logging" toml:"logging"`
	Debug   bool           `json:"debug" mapstructure:"debug" toml:"debug"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg)
}
