package apiclients

import (
	"net/http"

	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/metrics"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	counterName        metrics.CounterName = "api_clients"
	counterDescription string              = "number of API clients managed by the API client service"
	serviceName        string              = "api_clients_service"
)

var _ types.APIClientDataService = (*service)(nil)

type (
	// Config manages our body valdation.
	Config struct {
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
) types.APIClientDataService {
	return &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		cfg:                       cfg,
		apiClientDataManager:      clientDataManager,
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		encoderDecoder:            encoderDecoder,
		urlClientIDExtractor:      routeParamManager.BuildRouteParamStringIDFetcher(APIClientIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		secretGenerator:           secretGenerator,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}
}
