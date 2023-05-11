package tracing

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type errorHandler struct {
	logger logging.Logger
}

func (h errorHandler) Handle(err error) {
	h.logger.Error(err, "tracer reported issue")
}

func init() {
	// set this to a noop error handler just so one is set
	otel.SetErrorHandler(errorHandler{logger: zerolog.NewZerologLogger(logging.ErrorLevel).WithName("otel_errors")})
}

// Tracer describes a tracer.
type Tracer interface {
	StartSpan(ctx context.Context) (context.Context, Span)
	StartCustomSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span)
}

// TracerProvider is a simple alias for trace.TracerProvider.
type TracerProvider interface {
	trace.TracerProvider
	ForceFlush(ctx context.Context) error
}

type noopTracerProvider struct{}

func (n *noopTracerProvider) Tracer(instrumentationName string, opts ...trace.TracerOption) trace.Tracer {
	return trace.NewNoopTracerProvider().Tracer(instrumentationName, opts...)
}

func (n *noopTracerProvider) ForceFlush(_ context.Context) error {
	return nil
}

// NewNoopTracerProvider is a shadow for opentelemetry's NewNoopTracerProvider.
var NewNoopTracerProvider = func() TracerProvider {
	return &noopTracerProvider{}
}
