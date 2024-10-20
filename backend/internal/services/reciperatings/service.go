package reciperatings

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipe_ratings_service"
)

var _ types.RecipeRatingDataService = (*service)(nil)

type (
	// service handles valid instruments.
	service struct {
		logger                    logging.Logger
		recipeRatingDataManager   types.RecipeRatingDataManager
		recipeIDFetcher           func(*http.Request) string
		recipeRatingIDFetcher     func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new RecipeRatingsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	recipeRatingDataManager types.RecipeRatingDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.RecipeRatingDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeRatingIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(RecipeRatingIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeRatingDataManager:   recipeRatingDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
