package recipestepcompletionconditions

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
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipesteps"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipe_step_completion_conditions_service"
)

var _ types.RecipeStepCompletionConditionDataService = (*service)(nil)

type (
	// service handles recipe step ingredients.
	service struct {
		logger                                   logging.Logger
		recipeStepCompletionConditionDataManager types.RecipeStepCompletionConditionDataManager
		recipeIDFetcher                          func(*http.Request) string
		recipeStepIDFetcher                      func(*http.Request) string
		recipeStepCompletionConditionIDFetcher   func(*http.Request) string
		sessionContextDataFetcher                func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                     messagequeue.Publisher
		encoderDecoder                           encoding.ServerEncoderDecoder
		tracer                                   tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepCompletionConditionsService.
func ProvideService(
	logger logging.Logger,
	recipeStepCompletionConditionDataManager types.RecipeStepCompletionConditionDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.RecipeStepCompletionConditionDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                                   logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                          routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:                      routeParamManager.BuildRouteParamStringIDFetcher(recipestepsservice.RecipeStepIDURIParamKey),
		recipeStepCompletionConditionIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepCompletionConditionIDURIParamKey),
		sessionContextDataFetcher:                authentication.FetchContextFromRequest,
		recipeStepCompletionConditionDataManager: recipeStepCompletionConditionDataManager,
		dataChangesPublisher:                     dataChangesPublisher,
		encoderDecoder:                           encoder,
		tracer:                                   tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
