package webhooks

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
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
		sessionContextDataFetcher    func(*http.Request) (*sessions.ContextData, error)
	}
)

// ProvideWebhooksService builds a new WebhooksService.
func ProvideWebhooksService(
	logger logging.Logger,
	webhookDataManager types.WebhookDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.WebhookDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		webhookDataManager:           webhookDataManager,
		encoderDecoder:               encoder,
		dataChangesPublisher:         dataChangesPublisher,
		sessionContextDataFetcher:    sessions.FetchContextFromRequest,
		webhookIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(WebhookIDURIParamKey),
		webhookTriggerEventIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(WebhookTriggerEventIDURIParamKey),
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
