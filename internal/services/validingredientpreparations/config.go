package validingredientpreparations

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	Logging logging.Config `json:"logging" mapstructure:"logging" toml:"logging,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.Logging, validation.Required),
	)
}
