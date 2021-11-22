package apiclients

import (
	"net/http"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/random"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	counterName        metrics.CounterName = "api_clients"
	counterDescription string              = "number of API clients managed by the API client service"
	serviceName        string              = "api_clients_service"
)

var _ types.APIClientDataService = (*service)(nil)

type (
	config struct {
		minimumUsernameLength, minimumPasswordLength uint8
	}

	// service manages our API clients via HTTP.
	service struct {
		logger                    logging.Logger
		cfg                       *config
		apiClientDataManager      types.APIClientDataManager
		userDataManager           types.UserDataManager
		authenticator             authentication.Authenticator
		encoderDecoder            encoding.ServerEncoderDecoder
		urlClientIDExtractor      func(req *http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		apiClientCounter          metrics.UnitCounter
		secretGenerator           random.Generator
		tracer                    tracing.Tracer
		customerDataCollector     customerdata.Collector
	}
)

// ProvideAPIClientsService builds a new APIClientsService.
func ProvideAPIClientsService(
	logger logging.Logger,
	clientDataManager types.APIClientDataManager,
	userDataManager types.UserDataManager,
	authenticator authentication.Authenticator,
	encoderDecoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
	cfg *config,
	customerDataCollector customerdata.Collector,
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
		apiClientCounter:          metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		secretGenerator:           random.NewGenerator(logger),
		tracer:                    tracing.NewTracer(serviceName),
		customerDataCollector:     customerDataCollector,
	}
}
