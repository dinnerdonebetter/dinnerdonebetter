package recipestepinstruments

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	serviceName string = "recipe_step_instruments_service"
)

var _ types.RecipeStepInstrumentDataService = (*service)(nil)

type (
	// service handles recipe step instruments.
	service struct {
		logger                          logging.Logger
		recipeStepInstrumentDataManager types.RecipeStepInstrumentDataManager
		recipeIDFetcher                 func(*http.Request) string
		recipeStepIDFetcher             func(*http.Request) string
		recipeStepInstrumentIDFetcher   func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepInstrumentsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	recipeStepInstrumentDataManager types.RecipeStepInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.RecipeStepInstrumentDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step instruments service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                 routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(recipestepsservice.RecipeStepIDURIParamKey),
		recipeStepInstrumentIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepInstrumentIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		recipeStepInstrumentDataManager: recipeStepInstrumentDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
