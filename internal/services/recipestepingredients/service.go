package recipestepingredients

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	publishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	routing "github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/search"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "recipe_step_ingredients_service"
)

var _ types.RecipeStepIngredientDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipe step ingredients.
	service struct {
		logger                          logging.Logger
		recipeStepIngredientDataManager types.RecipeStepIngredientDataManager
		recipeIDFetcher                 func(*http.Request) string
		recipeStepIDFetcher             func(*http.Request) string
		recipeStepIngredientIDFetcher   func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher              publishers.Publisher
		preUpdatesPublisher             publishers.Publisher
		preArchivesPublisher            publishers.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepIngredientsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	recipeStepIngredientDataManager types.RecipeStepIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.RecipeStepIngredientDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step ingredient queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step ingredient queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step ingredient queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                 routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(recipestepsservice.RecipeStepIDURIParamKey),
		recipeStepIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepIngredientIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		recipeStepIngredientDataManager: recipeStepIngredientDataManager,
		preWritesPublisher:              preWritesPublisher,
		preUpdatesPublisher:             preUpdatesPublisher,
		preArchivesPublisher:            preArchivesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(serviceName),
	}

	return svc, nil
}
