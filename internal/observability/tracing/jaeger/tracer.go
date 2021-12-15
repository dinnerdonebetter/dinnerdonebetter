package jaeger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"

	"github.com/prixfixeco/api_server/internal/observability/logging"
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
func SetupJaeger(_ context.Context, c *Config) (trace.TracerProvider, error) {
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

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(c.ServiceName),
			// attribute.String(tagKey, tagVal),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
