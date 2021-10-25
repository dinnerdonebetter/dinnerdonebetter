package webhooks

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
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
		preWritesPublisher        publishers.Publisher
		preArchivesPublisher      publishers.Publisher
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
	publisherProvider publishers.PublisherProvider,
) (types.WebhookDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up pre-writes producer: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up pre-archives producer: %w", err)
	}

	s := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		webhookDataManager:        webhookDataManager,
		encoderDecoder:            encoder,
		preWritesPublisher:        preWritesPublisher,
		preArchivesPublisher:      preArchivesPublisher,
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		webhookIDFetcher:          routeParamManager.BuildRouteParamStringIDFetcher(WebhookIDURIParamKey),
		tracer:                    tracing.NewTracer(serviceName),
	}

	return s, nil
}
