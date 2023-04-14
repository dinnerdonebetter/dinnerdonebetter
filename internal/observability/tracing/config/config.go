package config

import (
	"context"
	"strings"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/observability/tracing/cloudtrace"
	"github.com/prixfixeco/backend/internal/observability/tracing/jaeger"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderJaeger represents the open source tracing server.
	ProviderJaeger = "jaeger"
	// ProviderCloudTrace represents the GCP Cloud Trace service.
	ProviderCloudTrace = "cloudtrace"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{}

		CloudTrace *cloudtrace.Config `json:"cloudTrace,omitempty" mapstructure:"cloud_trace" toml:"cloud_trace,omitempty"`
		Jaeger     *jaeger.Config     `json:"jaeger,omitempty" mapstructure:"jaeger" toml:"jaeger,omitempty"`
		Provider   string             `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
	}
)

// ProvideTracerProvider provides an instrumentation handler.
func (c *Config) ProvideTracerProvider(ctx context.Context, l logging.Logger) (traceProvider tracing.TracerProvider, err error) {
	logger := l.WithValue("tracing_provider", c.Provider)
	logger.Info("setting tracing provider")

	p := strings.TrimSpace(strings.ToLower(c.Provider))

	switch p {
	case ProviderJaeger:
		return jaeger.SetupJaeger(ctx, c.Jaeger)
	case ProviderCloudTrace:
		return cloudtrace.SetupCloudTrace(ctx, c.CloudTrace)
	default:
		logger.Debug("invalid tracing provider")
		return tracing.NewNoopTracerProvider(), nil
	}
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Provider, validation.In("", ProviderJaeger, ProviderCloudTrace)),
		validation.Field(&c.Jaeger, validation.When(c.Provider == ProviderJaeger, validation.Required).Else(validation.Nil)),
		validation.Field(&c.CloudTrace, validation.When(c.Provider == ProviderCloudTrace, validation.Required).Else(validation.Nil)),
	)
}
