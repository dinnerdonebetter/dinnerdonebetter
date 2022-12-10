package observability

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	loggingcfg "github.com/prixfixeco/backend/internal/observability/logging/config"
	metricscfg "github.com/prixfixeco/backend/internal/observability/metrics/config"
	tracingcfg "github.com/prixfixeco/backend/internal/observability/tracing/config"
)

type (
	// Config contains settings about how we report our metrics.
	Config struct {
		_ struct{}

		Logging loggingcfg.Config `json:"logging" mapstructure:"logging" toml:"logging,omitempty"`
		Metrics metricscfg.Config `json:"metrics" mapstructure:"metrics" toml:"metrics,omitempty"`
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
