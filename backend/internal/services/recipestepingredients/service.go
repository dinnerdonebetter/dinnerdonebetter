package recipestepingredients

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipe_step_ingredients_service"
)

var _ types.RecipeStepIngredientDataService = (*service)(nil)

type (
	// service handles recipe step ingredients.
	service struct {
		logger                          logging.Logger
		recipeStepIngredientDataManager types.RecipeStepIngredientDataManager
		recipeIDFetcher                 func(*http.Request) string
		recipeStepIDFetcher             func(*http.Request) string
		recipeStepIngredientIDFetcher   func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepIngredientsService.
func ProvideService(
	logger logging.Logger,
	recipeStepIngredientDataManager types.RecipeStepIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.RecipeStepIngredientDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                 routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(recipestepsservice.RecipeStepIDURIParamKey),
		recipeStepIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepIngredientIDURIParamKey),
		sessionContextDataFetcher:       authentication.FetchContextFromRequest,
		recipeStepIngredientDataManager: recipeStepIngredientDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
