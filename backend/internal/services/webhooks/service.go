package webhooks

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "webhooks_service"
)

var _ types.WebhookDataService = (*service)(nil)

type (
	// service handles webhooks.
	service struct {
		logger                       logging.Logger
		webhookDataManager           types.WebhookDataManager
		tracer                       tracing.Tracer
		encoderDecoder               encoding.ServerEncoderDecoder
		dataChangesPublisher         messagequeue.Publisher
		webhookIDFetcher             func(*http.Request) string
		webhookTriggerEventIDFetcher func(*http.Request) string
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
	}
)

// ProvideWebhooksService builds a new WebhooksService.
func ProvideWebhooksService(
	logger logging.Logger,
	cfg *Config,
	webhookDataManager types.WebhookDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.WebhookDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		webhookDataManager:           webhookDataManager,
		encoderDecoder:               encoder,
		dataChangesPublisher:         dataChangesPublisher,
		sessionContextDataFetcher:    authentication.FetchContextFromRequest,
		webhookIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(WebhookIDURIParamKey),
		webhookTriggerEventIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(WebhookTriggerEventIDURIParamKey),
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
