package pages

import (
	"net/url"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

type Builder struct {
	tracerProvider tracing.TracerProvider
	tracer         tracing.Tracer
	logger         logging.Logger
	apiServerURL   *url.URL
}

func NewPageBuilder(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	apiServerURL *url.URL,
) *Builder {
	s := &Builder{
		tracer:         tracing.NewTracer(tracerProvider.Tracer("admin-page-builder")),
		tracerProvider: tracerProvider,
		logger:         logger,
		apiServerURL:   apiServerURL,
	}

	return s
}
