package oauth2clients

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "oauth2_clients_service"
)

var _ types.OAuth2ClientDataService = (*service)(nil)

type (
	// service manages our OAuth2 clients via HTTP.
	service struct {
		logger                    logging.Logger
		oauth2ClientDataManager   types.OAuth2ClientDataManager
		encoderDecoder            encoding.ServerEncoderDecoder
		urlClientIDExtractor      func(req *http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		secretGenerator           random.Generator
		tracer                    tracing.Tracer
		dataChangesPublisher      messagequeue.Publisher
	}
)

// ProvideOAuth2ClientsService builds a new OAuth2ClientsService.
func ProvideOAuth2ClientsService(
	logger logging.Logger,
	clientDataManager types.OAuth2ClientDataManager,
	encoderDecoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	cfg *Config,
	tracerProvider tracing.TracerProvider,
	secretGenerator random.Generator,
	publisherProvider messagequeue.PublisherProvider,
) (types.OAuth2ClientDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		oauth2ClientDataManager:   clientDataManager,
		encoderDecoder:            encoderDecoder,
		urlClientIDExtractor:      routeParamManager.BuildRouteParamStringIDFetcher(OAuth2ClientIDURIParamKey),
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		secretGenerator:           secretGenerator,
		dataChangesPublisher:      dataChangesPublisher,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
