package validpreparations

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "valid_preparations"
	counterDescription string              = "the number of valid preparations managed by the valid preparations service"
	serviceName        string              = "valid_preparations_service"
)

var _ types.ValidPreparationDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid preparations.
	service struct {
		logger                      logging.Logger
		validPreparationDataManager types.ValidPreparationDataManager
		validPreparationIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		validPreparationCounter     metrics.UnitCounter
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
		search                      SearchIndex
	}
)

// ProvideService builds a new ValidPreparationsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	validPreparationDataManager types.ValidPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	searchIndexProvider search.IndexManagerProvider,
	routeParamManager routing.RouteParamManager,
) (types.ValidPreparationDataService, error) {
	searchIndexManager, err := searchIndexProvider(search.IndexPath(cfg.SearchIndexPath), "valid_preparations", logger)
	if err != nil {
		return nil, fmt.Errorf("setting up search index: %w", err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, ValidPreparationIDURIParamKey, "valid_preparation"),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		validPreparationDataManager: validPreparationDataManager,
		encoderDecoder:              encoder,
		validPreparationCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		search:                      searchIndexManager,
		tracer:                      tracing.NewTracer(serviceName),
	}

	return svc, nil
}
