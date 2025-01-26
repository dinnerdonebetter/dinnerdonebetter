package oauth2clients

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
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
		secretGenerator           random.Generator
		tracer                    tracing.Tracer
		dataChangesPublisher      messagequeue.Publisher
		urlClientIDExtractor      func(req *http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
	}
)

// ProvideOAuth2ClientsService builds a new OAuth2ClientsService.
func ProvideOAuth2ClientsService(
	logger logging.Logger,
	cfg *Config,
	clientDataManager types.OAuth2ClientDataManager,
	encoderDecoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	tracerProvider tracing.TracerProvider,
	secretGenerator random.Generator,
	publisherProvider messagequeue.PublisherProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.OAuth2ClientDataService, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError(serviceName)
	}

	if queueConfig == nil {
		return nil, internalerrors.NilConfigError(fmt.Sprintf("%s queues config", serviceName))
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
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
