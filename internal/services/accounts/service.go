package accounts

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "accounts"
	counterDescription string              = "the number of accounts managed by the accounts service"
	serviceName        string              = "accounts_service"
)

var _ types.AccountDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles to-do list accounts.
	service struct {
		logger                       logging.Logger
		accountDataManager           types.AccountDataManager
		accountMembershipDataManager types.AccountUserMembershipDataManager
		tracer                       tracing.Tracer
		accountCounter               metrics.UnitCounter
		encoderDecoder               encoding.ServerEncoderDecoder
		preWritesPublisher           publishers.Publisher
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
		userIDFetcher                func(*http.Request) string
		accountIDFetcher             func(*http.Request) string
	}
)

// ProvideService builds a new AccountsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	accountDataManager types.AccountDataManager,
	accountMembershipDataManager types.AccountUserMembershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.AccountDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up event publisher: %w", err)
	}

	s := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		accountIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(AccountIDURIParamKey),
		userIDFetcher:                routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher:    authservice.FetchContextFromRequest,
		accountDataManager:           accountDataManager,
		accountMembershipDataManager: accountMembershipDataManager,
		encoderDecoder:               encoder,
		preWritesPublisher:           preWritesPublisher,
		accountCounter:               metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                       tracing.NewTracer(serviceName),
	}

	return s, nil
}
