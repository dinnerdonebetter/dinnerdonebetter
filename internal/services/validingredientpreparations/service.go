package validingredientpreparations

import (
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
	counterName        metrics.CounterName = "valid_ingredient_preparations"
	counterDescription string              = "the number of valid ingredient preparations managed by the valid ingredient preparations service"
	serviceName        string              = "valid_ingredient_preparations_service"
)

var _ types.ValidIngredientPreparationDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid ingredient preparations.
	service struct {
		logger                                logging.Logger
		validIngredientPreparationDataManager types.ValidIngredientPreparationDataManager
		validIngredientPreparationIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
		validIngredientPreparationCounter     metrics.UnitCounter
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientPreparationsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	validIngredientPreparationDataManager types.ValidIngredientPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.ValidIngredientPreparationDataService, error) {
	svc := &service{
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientPreparationIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, ValidIngredientPreparationIDURIParamKey, "valid_ingredient_preparation"),
		sessionContextDataFetcher:             authservice.FetchContextFromRequest,
		validIngredientPreparationDataManager: validIngredientPreparationDataManager,
		encoderDecoder:                        encoder,
		validIngredientPreparationCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                                tracing.NewTracer(serviceName),
	}

	return svc, nil
}
