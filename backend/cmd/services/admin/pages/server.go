package pages

import (
	"net/url"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
)

type PageBuilder struct {
	tracerProvider tracing.TracerProvider
	tracer         tracing.Tracer
	logger         logging.Logger
	router         routing.Router
	apiServerURL   *url.URL
}

func NewPageBuilder(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	router routing.Router,
	apiServerURL *url.URL,
) *PageBuilder {
	tracer := tracing.NewTracer(tracerProvider.Tracer("admin-page-builder"))

	s := &PageBuilder{
		tracer:         tracer,
		tracerProvider: tracerProvider,
		logger:         logger,
		router:         router,
		apiServerURL:   apiServerURL,
	}

	return s
}
