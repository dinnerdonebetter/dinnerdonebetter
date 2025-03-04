package recipemanagement

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/images"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipes_service"
)

var _ types.RecipeManagementDataService = (*service)(nil)

type (
	// service handles recipes.
	service struct {
		logger                                 logging.Logger
		tracer                                 tracing.Tracer
		recipeManagementDataManager            types.RecipeManagementDataManager
		recipePrepTaskIDFetcher                func(*http.Request) string
		recipeRatingIDFetcher                  func(*http.Request) string
		recipeIDFetcher                        func(*http.Request) string
		recipeStepIDFetcher                    func(*http.Request) string
		recipeStepVesselIDFetcher              func(*http.Request) string
		recipeStepProductIDFetcher             func(*http.Request) string
		recipeStepInstrumentIDFetcher          func(*http.Request) string
		recipeStepIngredientIDFetcher          func(*http.Request) string
		recipeStepCompletionConditionIDFetcher func(*http.Request) string
		// TODO: recipeAnalyzer                         recipeanalysis.RecipeAnalyzer
		imageUploadProcessor      images.MediaUploadProcessor
		encoderDecoder            encoding.ServerEncoderDecoder
		dataChangesPublisher      messagequeue.Publisher
		searchIndex               textsearch.IndexSearcher[eatingindexing.RecipeSearchSubset]
		uploadManager             uploads.UploadManager
		sessionContextDataFetcher func(*http.Request) (*sessions.ContextData, error)
		cfg                       *Config
	}
)

var errInvalidConfig = errors.New("config cannot be nil")

// ProvideService builds a new RecipesService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	// TODO: recipeGrapher recipeanalysis.RecipeAnalyzer,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	imageUploadProcessor images.MediaUploadProcessor,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	queueConfig *msgconfig.QueuesConfig,
	recipesDataManager types.RecipeManagementDataManager,
) (types.RecipeManagementDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	if cfg == nil {
		return nil, errInvalidConfig
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	uploader, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Uploads.Storage, routeParamManager)
	if err != nil {
		return nil, fmt.Errorf("initializing %s upload manager: %w", serviceName, err)
	}

	searchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.RecipeSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeRecipes)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing recipe index manager")
	}

	svc := &service{
		logger:                                 logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:                        routeParamManager.BuildRouteParamStringIDFetcher(RecipeIDURIParamKey),
		recipeStepIDFetcher:                    routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepIDURIParamKey),
		recipeStepVesselIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepVesselIDURIParamKey),
		recipeStepProductIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepProductIDURIParamKey),
		recipeStepInstrumentIDFetcher:          routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepInstrumentIDURIParamKey),
		recipeStepIngredientIDFetcher:          routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepIngredientIDURIParamKey),
		recipeStepCompletionConditionIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepCompletionConditionIDURIParamKey),
		recipePrepTaskIDFetcher:                routeParamManager.BuildRouteParamStringIDFetcher(RecipePrepTaskIDURIParamKey),
		recipeRatingIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(RecipeRatingIDURIParamKey),
		sessionContextDataFetcher:              sessions.FetchContextFromRequest,
		recipeManagementDataManager:            recipesDataManager,
		cfg:                                    cfg,
		dataChangesPublisher:                   dataChangesPublisher,
		encoderDecoder:                         encoder,
		// TODO: recipeAnalyzer:                         recipeGrapher,
		uploadManager:        uploader,
		imageUploadProcessor: imageUploadProcessor,
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:          searchIndex,
	}

	return svc, nil
}
