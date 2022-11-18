package config

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/backend/internal/customerdata"
	"github.com/prixfixeco/backend/internal/customerdata/segment"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

const (
	// ProviderSegment represents Segment.
	ProviderSegment = "segment"
)

type (
	// Config is the configuration structure.
	Config struct {
		Provider string `json:"provider" mapstructure:"provider" toml:"provider,omitempty"`
		APIToken string `json:"apiToken" mapstructure:"api_token" toml:"api_token,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.APIToken, validation.When(strings.EqualFold(strings.TrimSpace(cfg.Provider), ProviderSegment), validation.Required)),
	)
}

// ProvideCollector provides a collector.
func (cfg *Config) ProvideCollector(logger logging.Logger, tracerProvider tracing.TracerProvider) (customerdata.Collector, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderSegment:
		return segment.NewSegmentCustomerDataCollector(logger, tracerProvider, cfg.APIToken)
	default:
		return customerdata.NewNoopCollector(), nil
	}
}
