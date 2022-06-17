package webhooks

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "webhooks_service"
)

var (
	_ types.WebhookDataService = (*service)(nil)
)

type (
	// service handles webhooks.
	service struct {
		logger                    logging.Logger
		webhookDataManager        types.WebhookDataManager
		tracer                    tracing.Tracer
		encoderDecoder            encoding.ServerEncoderDecoder
		dataChangesPublisher      messagequeue.Publisher
		webhookIDFetcher          func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
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
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up data changes publisher: %w", err)
	}

	s := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		webhookDataManager:        webhookDataManager,
		encoderDecoder:            encoder,
		dataChangesPublisher:      dataChangesPublisher,
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		webhookIDFetcher:          routeParamManager.BuildRouteParamStringIDFetcher(WebhookIDURIParamKey),
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return s, nil
}
