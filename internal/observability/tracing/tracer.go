package tracing

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

type errorHandler struct {
	logger logging.Logger
}

func (h errorHandler) Handle(err error) {
	h.logger.Error(err, "tracer reported issue")
}

func init() {
	// set this to a noop error handler just so one is set
	otel.SetErrorHandler(errorHandler{logger: zerolog.NewZerologLogger().WithName("otel_errors")})
}

// Tracer describes a tracer.
type Tracer interface {
	StartSpan(ctx context.Context) (context.Context, Span)
	StartCustomSpan(ctx context.Context, name string) (context.Context, Span)
}

// TracerProvider is a simple alias for trace.TracerProvider.
type TracerProvider trace.TracerProvider
