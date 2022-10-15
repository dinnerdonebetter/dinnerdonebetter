package recipepreptasks

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "recipe_prep_tasks_service"
)

var _ types.RecipePrepTaskDataService = (*service)(nil)

type (
	// service handles recipe prep tasks.
	service struct {
		logger                    logging.Logger
		recipePrepTaskDataManager types.RecipePrepTaskDataManager
		recipeIDFetcher           func(*http.Request) string
		recipePrepTaskIDFetcher   func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new RecipePrepTasksService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	recipePrepTaskDataManager types.RecipePrepTaskDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.RecipePrepTaskDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe prep tasks service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipePrepTaskIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipePrepTaskIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipePrepTaskDataManager: recipePrepTaskDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
