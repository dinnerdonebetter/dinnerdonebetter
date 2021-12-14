package config

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/trace"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing/jaeger"
)

const (
	// Jaeger represents the open source tracing server.
	Jaeger = "jaeger"
	// XRay represents the AWS tracing server.
	XRay = "xray"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{}

		Jaeger   *jaeger.Config `json:"jaeger,omitempty" mapstructure:"jaeger" toml:"jaeger,omitempty"`
		Provider string         `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
	}
)

// Initialize provides an instrumentation handler.
func (c *Config) Initialize(l logging.Logger) (traceProvidier trace.TracerProvider, flushFunc func(), err error) {
	logger := l.WithValue("tracing_provider", c.Provider)
	logger.Info("setting tracing provider")

	p := strings.TrimSpace(strings.ToLower(c.Provider))

	switch p {
	case Jaeger:
		logger.Debug("setting up jaeger")
		return jaeger.SetupJaeger(c.Jaeger)
	case XRay:
		return trace.NewNoopTracerProvider(), func() {}, nil
	default:
		logger.Debug("invalid tracing config")
		return nil, nil, fmt.Errorf("invalid tracing provider: %q", p)
	}
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Provider, validation.In("", Jaeger)),
		validation.Field(&c.Jaeger, validation.When(c.Provider == Jaeger, validation.Required).Else(validation.Nil)),
	)
}
