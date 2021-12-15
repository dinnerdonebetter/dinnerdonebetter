package observability

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability/metrics/config"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	tracingcfg "github.com/prixfixeco/api_server/internal/observability/tracing/config"
)

type (
	// Config contains settings about how we report our metrics.
	Config struct {
		_       struct{}
		Metrics config.Config     `json:"metrics" mapstructure:"metrics" toml:"metrics,omitempty"`
		Tracing tracingcfg.Config `json:"tracing" mapstructure:"tracing" toml:"tracing,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Metrics),
		validation.Field(&cfg.Tracing),
	)
}
