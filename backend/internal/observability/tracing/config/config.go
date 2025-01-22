package tracingcfg

import (
	"context"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/cloudtrace"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/oteltrace"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderOtel represents the open source tracing server.
	ProviderOtel = "otelgrpc"
	// ProviderCloudTrace represents the GCP Cloud Trace service.
	ProviderCloudTrace = "cloudtrace"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{} `json:"-"`

		CloudTrace                *cloudtrace.Config `env:"init"                                envPrefix:"CLOUDTRACE_"                    json:"cloudTrace,omitempty"`
		Otel                      *oteltrace.Config  `env:"init"                                envPrefix:"OTELGRPC_"                      json:"otelgrpc,omitempty"`
		ServiceName               string             `env:"TRACING_SERVICE_NAME"                json:"service_name,omitempty"`
		Provider                  string             `env:"TRACING_PROVIDER"                    json:"provider,omitempty"`
		SpanCollectionProbability float64            `env:"TRACING_SPAN_COLLECTION_PROBABILITY" json:"spanCollectionProbability,omitempty"`
	}
)

// ProvideTracerProvider provides a TracerProvider.
func (c *Config) ProvideTracerProvider(ctx context.Context, l logging.Logger) (tracing.TracerProvider, error) {
	logger := l.WithValue("tracing_provider", c.Provider)

	p := strings.TrimSpace(strings.ToLower(c.Provider))

	switch p {
	case ProviderOtel:
		logger.WithValue("otel", c.Otel).Info("configuring otelgrpc provider")
		tp, err := oteltrace.SetupOtelGRPC(ctx, c.ServiceName, c.SpanCollectionProbability, c.Otel)
		if err != nil {
			return nil, fmt.Errorf("configuring otelgrpc provider: %w", err)
		}

		return tp, nil
	case ProviderCloudTrace:
		logger.Info("configuring cloud trace provider")
		tp, err := cloudtrace.SetupCloudTrace(ctx, c.ServiceName, c.SpanCollectionProbability, c.CloudTrace)
		if err != nil {
			return nil, fmt.Errorf("configuring cloud trace provider: %w", err)
		}

		return tp, nil
	default:
		logger.Info("invalid tracing provider")
		return tracing.NewNoopTracerProvider(), nil
	}
}

// ProvideTracer provides an instrumentation handler.
func (c *Config) ProvideTracer(ctx context.Context, l logging.Logger, name string) (tracing.Tracer, error) {
	tp, err := c.ProvideTracerProvider(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("configuring tracing provider: %w", err)
	}

	return tracing.NewTracer(tp.Tracer(name)), nil
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Provider, validation.In("", ProviderOtel, ProviderCloudTrace)),
		validation.Field(&c.Otel, validation.When(c.Provider == ProviderOtel, validation.Required).Else(validation.Nil)),
		validation.Field(&c.CloudTrace, validation.When(c.Provider == ProviderCloudTrace, validation.Required).Else(validation.Nil)),
		validation.Field(&c.ServiceName, validation.Required),
		validation.Field(&c.SpanCollectionProbability, validation.Required),
	)
}
