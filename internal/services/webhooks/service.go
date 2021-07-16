package webhooks

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "webhooks"
	counterDescription string              = "the number of webhooks managed by the webhooks service"
	serviceName        string              = "webhooks_service"
)

var (
	_ types.WebhookDataService = (*service)(nil)
)

type (
	// service handles webhooks.
	service struct {
		logger                    logging.Logger
		webhookCounter            metrics.UnitCounter
		webhookDataManager        types.WebhookDataManager
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		webhookIDFetcher          func(*http.Request) uint64
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideWebhooksService builds a new WebhooksService.
func ProvideWebhooksService(
	logger logging.Logger,
	webhookDataManager types.WebhookDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) types.WebhookDataService {
	return &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		webhookDataManager:        webhookDataManager,
		encoderDecoder:            encoder,
		webhookCounter:            metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		webhookIDFetcher:          routeParamManager.BuildRouteParamIDFetcher(logger, WebhookIDURIParamKey, "webhook"),
		tracer:                    tracing.NewTracer(serviceName),
	}
}
