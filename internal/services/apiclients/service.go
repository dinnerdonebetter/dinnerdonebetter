package apiclients

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/random"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
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
		urlClientIDExtractor      func(req *http.Request) uint64
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		apiClientCounter          metrics.UnitCounter
		secretGenerator           random.Generator
		tracer                    tracing.Tracer
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
) types.APIClientDataService {
	return &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		cfg:                       cfg,
		apiClientDataManager:      clientDataManager,
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		encoderDecoder:            encoderDecoder,
		urlClientIDExtractor:      routeParamManager.BuildRouteParamIDFetcher(logger, APIClientIDURIParamKey, "api client"),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		apiClientCounter:          metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		secretGenerator:           random.NewGenerator(logger),
		tracer:                    tracing.NewTracer(serviceName),
	}
}
