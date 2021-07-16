package recipesteps

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "recipe_steps"
	counterDescription string              = "the number of recipe steps managed by the recipe steps service"
	serviceName        string              = "recipe_steps_service"
)

var _ types.RecipeStepDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipe steps.
	service struct {
		logger                    logging.Logger
		recipeStepDataManager     types.RecipeStepDataManager
		recipeIDFetcher           func(*http.Request) uint64
		recipeStepIDFetcher       func(*http.Request) uint64
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		recipeStepCounter         metrics.UnitCounter
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	recipeStepDataManager types.RecipeStepDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.RecipeStepDataService, error) {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamIDFetcher(logger, recipesservice.RecipeIDURIParamKey, "recipe"),
		recipeStepIDFetcher:       routeParamManager.BuildRouteParamIDFetcher(logger, RecipeStepIDURIParamKey, "recipe_step"),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeStepDataManager:     recipeStepDataManager,
		encoderDecoder:            encoder,
		recipeStepCounter:         metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                    tracing.NewTracer(serviceName),
	}

	return svc, nil
}
