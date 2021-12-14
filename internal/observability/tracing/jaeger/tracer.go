package jaeger

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

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
func SetupJaeger(c *Config) (trace.TracerProvider, func(), error) {
	// Create and install Jaeger export pipeline.
	tp, flush, err := jaeger.NewExportPipeline(
		jaeger.WithCollectorEndpoint(c.CollectorEndpoint),
		jaeger.WithProcessFromEnv(),
		jaeger.WithSDKOptions(
			sdktrace.WithSampler(sdktrace.TraceIDRatioBased(c.SpanCollectionProbability)),
			sdktrace.WithResource(resource.NewWithAttributes(
				attribute.String("exporter", "jaeger"),
				attribute.String("service.name", c.ServiceName),
			)),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing Jaeger: %w", err)
	}

	otel.SetTracerProvider(tp)

	return tp, flush, nil
}
