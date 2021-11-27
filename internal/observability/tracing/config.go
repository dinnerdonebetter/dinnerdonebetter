package tracing

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

const (
	// Jaeger represents the popular distributed tracing server.
	Jaeger = "jaeger"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{}

		Jaeger                    *JaegerConfig `json:"jaeger,omitempty" mapstructure:"jaeger" toml:"jaeger,omitempty"`
		Provider                  string        `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
		SpanCollectionProbability float64       `json:"spanCollectionProbability,omitempty" mapstructure:"span_collection_probability" toml:"span_collection_probability,omitempty"`
	}

	// JaegerConfig contains settings related to tracing with Jaeger.
	JaegerConfig struct {
		_ struct{}

		CollectorEndpoint string `json:"collector_endpoint,omitempty" mapstructure:"collector_endpoint" toml:"collector_endpoint,omitempty"`
		ServiceName       string `json:"service_name,omitempty" mapstructure:"service_name" toml:"service_name,omitempty"`
	}
)

// Initialize provides an instrumentation handler.
func (c *Config) Initialize(l logging.Logger) (flushFunc func(), err error) {
	logger := l.WithValue("tracing_provider", c.Provider)
	logger.Info("setting tracing provider")

	p := strings.TrimSpace(strings.ToLower(c.Provider))

	switch p {
	case Jaeger:
		logger.Debug("setting up jaeger")
		return c.SetupJaeger()
	case "":
		return nil, nil
	default:
		logger.Debug("invalid tracing config")
		return nil, fmt.Errorf("invalid tracing provider: %q", p)
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

var _ validation.ValidatableWithContext = (*JaegerConfig)(nil)

// ValidateWithContext validates the config struct.
func (c *JaegerConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.CollectorEndpoint, validation.Required),
		validation.Field(&c.ServiceName, validation.Required),
	)
}
