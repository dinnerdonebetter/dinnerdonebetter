package recipesteps

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	publishers "gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName string = "recipe_steps_service"
)

var _ types.RecipeStepDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipe steps.
	service struct {
		logger                    logging.Logger
		recipeStepDataManager     types.RecipeStepDataManager
		recipeIDFetcher           func(*http.Request) string
		recipeStepIDFetcher       func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher        publishers.Publisher
		preUpdatesPublisher       publishers.Publisher
		preArchivesPublisher      publishers.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	recipeStepDataManager types.RecipeStepDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.RecipeStepDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:       routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeStepDataManager:     recipeStepDataManager,
		preWritesPublisher:        preWritesPublisher,
		preUpdatesPublisher:       preUpdatesPublisher,
		preArchivesPublisher:      preArchivesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(serviceName),
	}

	return svc, nil
}
