package recipes

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
	counterName        metrics.CounterName = "recipes"
	counterDescription string              = "the number of recipes managed by the recipes service"
	serviceName        string              = "recipes_service"
)

var _ types.RecipeDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipes.
	service struct {
		logger                    logging.Logger
		recipeDataManager         types.RecipeDataManager
		recipeIDFetcher           func(*http.Request) uint64
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		recipeCounter             metrics.UnitCounter
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new RecipesService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	recipeDataManager types.RecipeDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.RecipeDataService, error) {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamIDFetcher(logger, RecipeIDURIParamKey, "recipe"),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeDataManager:         recipeDataManager,
		encoderDecoder:            encoder,
		recipeCounter:             metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                    tracing.NewTracer(serviceName),
	}

	return svc, nil
}
