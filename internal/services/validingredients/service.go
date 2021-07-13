package validingredients

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
	counterName        metrics.CounterName = "valid_ingredients"
	counterDescription string              = "the number of valid ingredients managed by the valid ingredients service"
	serviceName        string              = "valid_ingredients_service"
)

var _ types.ValidIngredientDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid ingredients.
	service struct {
		logger                     logging.Logger
		validIngredientDataManager types.ValidIngredientDataManager
		validIngredientIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		validIngredientCounter     metrics.UnitCounter
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
		search                     SearchIndex
	}
)

// ProvideService builds a new ValidIngredientsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	validIngredientDataManager types.ValidIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	searchIndexProvider search.IndexManagerProvider,
	routeParamManager routing.RouteParamManager,
) (types.ValidIngredientDataService, error) {
	searchIndexManager, err := searchIndexProvider(search.IndexPath(cfg.SearchIndexPath), "valid_ingredients", logger)
	if err != nil {
		return nil, fmt.Errorf("setting up search index: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, ValidIngredientIDURIParamKey, "valid_ingredient"),
		sessionContextDataFetcher:  authservice.FetchContextFromRequest,
		validIngredientDataManager: validIngredientDataManager,
		encoderDecoder:             encoder,
		validIngredientCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		search:                     searchIndexManager,
		tracer:                     tracing.NewTracer(serviceName),
	}

	return svc, nil
}
