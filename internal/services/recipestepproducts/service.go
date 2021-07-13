package recipestepproducts

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
	counterName        metrics.CounterName = "recipe_step_products"
	counterDescription string              = "the number of recipe step products managed by the recipe step products service"
	serviceName        string              = "recipe_step_products_service"
)

var _ types.RecipeStepProductDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipe step products.
	service struct {
		logger                       logging.Logger
		recipeStepProductDataManager types.RecipeStepProductDataManager
		recipeIDFetcher              func(*http.Request) uint64
		recipeStepIDFetcher          func(*http.Request) uint64
		recipeStepProductIDFetcher   func(*http.Request) uint64
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
		recipeStepProductCounter     metrics.UnitCounter
		encoderDecoder               encoding.ServerEncoderDecoder
		tracer                       tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepProductsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	recipeStepProductDataManager types.RecipeStepProductDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.RecipeStepProductDataService, error) {
	svc := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:              routeParamManager.BuildRouteParamIDFetcher(logger, recipesservice.RecipeIDURIParamKey, "recipe"),
		recipeStepIDFetcher:          routeParamManager.BuildRouteParamIDFetcher(logger, recipestepsservice.RecipeStepIDURIParamKey, "recipe_step"),
		recipeStepProductIDFetcher:   routeParamManager.BuildRouteParamIDFetcher(logger, RecipeStepProductIDURIParamKey, "recipe_step_product"),
		sessionContextDataFetcher:    authservice.FetchContextFromRequest,
		recipeStepProductDataManager: recipeStepProductDataManager,
		encoderDecoder:               encoder,
		recipeStepProductCounter:     metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                       tracing.NewTracer(serviceName),
	}

	return svc, nil
}
