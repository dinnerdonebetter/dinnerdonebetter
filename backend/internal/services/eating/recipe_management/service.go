package recipemanagement

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	textsearch "github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
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
		recipeAnalyzer                         recipeanalysis.RecipeAnalyzer
		imageUploadProcessor                   images.MediaUploadProcessor
		encoderDecoder                         encoding.ServerEncoderDecoder
		dataChangesPublisher                   messagequeue.Publisher
		searchIndex                            textsearch.IndexSearcher[types.RecipeSearchSubset]
		uploadManager                          uploads.UploadManager
		sessionContextDataFetcher              func(*http.Request) (*types.SessionContextData, error)
		cfg                                    *Config
	}
)

var errInvalidConfig = errors.New("config cannot be nil")

// ProvideService builds a new RecipesService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	recipeGrapher recipeanalysis.RecipeAnalyzer,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	imageUploadProcessor images.MediaUploadProcessor,
	tracerProvider tracing.TracerProvider,
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

	searchIndex, err := textsearchcfg.ProvideIndex[types.RecipeSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeRecipes)
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
		sessionContextDataFetcher:              authentication.FetchContextFromRequest,
		recipeManagementDataManager:            recipesDataManager,
		cfg:                                    cfg,
		dataChangesPublisher:                   dataChangesPublisher,
		encoderDecoder:                         encoder,
		recipeAnalyzer:                         recipeGrapher,
		uploadManager:                          uploader,
		imageUploadProcessor:                   imageUploadProcessor,
		tracer:                                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:                            searchIndex,
	}

	return svc, nil
}
