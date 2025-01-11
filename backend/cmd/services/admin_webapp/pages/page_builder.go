package pages

import (
	"net/url"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
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
	tracer := tracing.NewTracer(tracerProvider.Tracer("admin-page-builder"))

	s := &PageBuilder{
		tracer:         tracer,
		tracerProvider: tracerProvider,
		logger:         logger,
		apiServerURL:   apiServerURL,
	}

	return s
}

func (b *PageBuilder) buildAPIClient() (*apiclient.Client, error) {
	return apiclient.NewClient(b.apiServerURL, b.tracerProvider)
}
