package apiclients

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	serviceName string = "api_clients_service"
)

var _ types.APIClientDataService = (*service)(nil)

type (
	// Config manages our body validation.
	Config struct {
		dataChangesTopicName string
		minimumUsernameLength,
		minimumPasswordLength uint8
	}

	// service manages our API clients via HTTP.
	service struct {
		logger                    logging.Logger
		cfg                       *Config
		apiClientDataManager      types.APIClientDataManager
		userDataManager           types.UserDataManager
		authenticator             authentication.Authenticator
		encoderDecoder            encoding.ServerEncoderDecoder
		urlClientIDExtractor      func(req *http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		secretGenerator           random.Generator
		tracer                    tracing.Tracer
		dataChangesPublisher      messagequeue.Publisher
	}
)

// ProvideAPIClientsService builds a new APIClientsService.
func ProvideAPIClientsService(
	logger logging.Logger,
	clientDataManager types.APIClientDataManager,
	userDataManager types.UserDataManager,
	authenticator authentication.Authenticator,
	encoderDecoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	cfg *Config,
	tracerProvider tracing.TracerProvider,
	secretGenerator random.Generator,
	publisherProvider messagequeue.PublisherProvider,
) (types.APIClientDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.dataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up api clients service data changes publisher: %w", err)
	}

	s := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		cfg:                       cfg,
		apiClientDataManager:      clientDataManager,
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		encoderDecoder:            encoderDecoder,
		urlClientIDExtractor:      routeParamManager.BuildRouteParamStringIDFetcher(APIClientIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		secretGenerator:           secretGenerator,
		dataChangesPublisher:      dataChangesPublisher,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return s, nil
}
