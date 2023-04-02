package jaeger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type errorHandler struct {
	logger logging.Logger
}

func (h errorHandler) Handle(err error) {
	h.logger.Error(err, "tracer reported issue")
}

func init() {
	otel.SetErrorHandler(errorHandler{logger: logging.NewNoopLogger().WithName("otel_errors")})
}

// SetupJaeger creates a new trace provider instance and registers it as global trace provider.
func SetupJaeger(ctx context.Context, c *Config) (tracing.TracerProvider, error) {
	// Create and install Jaeger export pipeline.
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(c.CollectorEndpoint),
			jaeger.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("initializing Jaeger: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithAttributes(semconv.ServiceNameKey.String(c.ServiceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("setting up process runtime version: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(c.SpanCollectionProbability)),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
