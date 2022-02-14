package config

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/observability/tracing/cloudtrace"
	"github.com/prixfixeco/api_server/internal/observability/tracing/jaeger"
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

// Initialize provides an instrumentation handler.
func (c *Config) Initialize(ctx context.Context, l logging.Logger) (traceProvider tracing.TracerProvider, err error) {
	logger := l.WithValue("tracing_provider", c.Provider)
	logger.Info("setting tracing provider")

	p := strings.TrimSpace(strings.ToLower(c.Provider))

	switch p {
	case ProviderJaeger:
		return jaeger.SetupJaeger(ctx, c.Jaeger)
	case ProviderCloudTrace:
		return cloudtrace.SetupCloudTrace(ctx, c.CloudTrace)
	case "":
		return tracing.NewNoopTracerProvider(), nil
	default:
		logger.Debug("invalid tracing provider")
		return nil, fmt.Errorf("invalid tracing provider: %q", p)
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
