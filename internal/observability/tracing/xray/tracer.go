package xray

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

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

// SetupXRay creates a new trace provider instance and registers it as global trace provider.
func SetupXRay(ctx context.Context, c *Config) (trace.TracerProvider, error) {
	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return func() (trace.TracerProvider, error) {
		traceExporter, err := otlptracegrpc.New(
			grpcCtx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(c.CollectorEndpoint),
			otlptracegrpc.WithDialOption(
				grpc.WithBlock(),
			),
		)
		if err != nil {
			return nil, err
		}

		res := resource.NewWithAttributes(
			semconv.SchemaURL,
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(c.ServiceName),
		)

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.TraceIDRatioBased(c.SpanCollectionProbability)),
			sdktrace.WithResource(res),
			sdktrace.WithBatcher(traceExporter),
			sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		)

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(xray.Propagator{})

		return tp, nil
	}()
}
