package recipes

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/graphing"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "recipes_service"
)

var _ types.RecipeDataService = (*service)(nil)

type (
	// service handles recipes.
	service struct {
		logger                    logging.Logger
		recipeDataManager         types.RecipeDataManager
		recipeGrapher             graphing.RecipeDAGDiagramGenerator
		recipeIDFetcher           func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new RecipesService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	recipeDataManager types.RecipeDataManager,
	recipeGrapher graphing.RecipeDAGDiagramGenerator,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.RecipeDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(RecipeIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeDataManager:         recipeDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		recipeGrapher:             recipeGrapher,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
