package pages

import (
	"net/url"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

type CookieBuilder interface {
	Encode(name string, value any) (string, error)
}

type PageBuilder struct {
	tracerProvider tracing.TracerProvider
	tracer         tracing.Tracer
	logger         logging.Logger
	apiServerURL   *url.URL
}

func NewPageBuilder(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	apiServerURL *url.URL,
) *PageBuilder {
	s := &PageBuilder{
		tracer:         tracing.NewTracer(tracerProvider.Tracer("admin-page-builder")),
		tracerProvider: tracerProvider,
		logger:         logger,
		apiServerURL:   apiServerURL,
	}

	return s
}
