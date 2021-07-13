package recipestepingredients

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
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "recipe_step_ingredients"
	counterDescription string              = "the number of recipe step ingredients managed by the recipe step ingredients service"
	serviceName        string              = "recipe_step_ingredients_service"
)

var _ types.RecipeStepIngredientDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipe step ingredients.
	service struct {
		logger                          logging.Logger
		recipeStepIngredientDataManager types.RecipeStepIngredientDataManager
		recipeIDFetcher                 func(*http.Request) uint64
		recipeStepIDFetcher             func(*http.Request) uint64
		recipeStepIngredientIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		recipeStepIngredientCounter     metrics.UnitCounter
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepIngredientsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	recipeStepIngredientDataManager types.RecipeStepIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.RecipeStepIngredientDataService, error) {
	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                 routeParamManager.BuildRouteParamIDFetcher(logger, recipesservice.RecipeIDURIParamKey, "recipe"),
		recipeStepIDFetcher:             routeParamManager.BuildRouteParamIDFetcher(logger, recipestepsservice.RecipeStepIDURIParamKey, "recipe_step"),
		recipeStepIngredientIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, RecipeStepIngredientIDURIParamKey, "recipe_step_ingredient"),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		recipeStepIngredientDataManager: recipeStepIngredientDataManager,
		encoderDecoder:                  encoder,
		recipeStepIngredientCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                          tracing.NewTracer(serviceName),
	}

	return svc, nil
}
