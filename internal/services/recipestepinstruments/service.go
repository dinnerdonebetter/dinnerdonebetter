package recipestepinstruments

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
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName string = "recipe_step_instruments_service"
)

var _ types.RecipeStepInstrumentDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles recipe step instruments.
	service struct {
		logger                          logging.Logger
		recipeStepInstrumentDataManager types.RecipeStepInstrumentDataManager
		recipeIDFetcher                 func(*http.Request) string
		recipeStepIDFetcher             func(*http.Request) string
		recipeStepInstrumentIDFetcher   func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher              publishers.Publisher
		preUpdatesPublisher             publishers.Publisher
		preArchivesPublisher            publishers.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepInstrumentsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	recipeStepInstrumentDataManager types.RecipeStepInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.RecipeStepInstrumentDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step instrument queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step instrument queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step instrument queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                 routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(recipestepsservice.RecipeStepIDURIParamKey),
		recipeStepInstrumentIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepInstrumentIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		recipeStepInstrumentDataManager: recipeStepInstrumentDataManager,
		preWritesPublisher:              preWritesPublisher,
		preUpdatesPublisher:             preUpdatesPublisher,
		preArchivesPublisher:            preArchivesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(serviceName),
	}

	return svc, nil
}
