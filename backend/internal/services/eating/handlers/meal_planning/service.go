package mealplanning

import (
	"context"
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
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "meals_service"
)

var _ types.MealDataService = (*service)(nil)

type (
	// service handles meals.
	service struct {
		logger                                logging.Logger
		dataChangesPublisher                  messagequeue.Publisher
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
		mealPlanningDataManager               types.MealPlanningDataManager
		mealsSearchIndex                      textsearch.IndexSearcher[eatingindexing.MealSearchSubset]
		householdInstrumentOwnershipIDFetcher func(*http.Request) string
		userIngredientPreferenceIDFetcher     func(*http.Request) string
		mealPlanIDFetcher                     func(*http.Request) string
		mealPlanGroceryListItemIDFetcher      func(*http.Request) string
		mealPlanEventIDFetcher                func(*http.Request) string
		mealIDFetcher                         func(*http.Request) string
		mealPlanOptionIDFetcher               func(*http.Request) string
		mealPlanTaskIDFetcher                 func(*http.Request) string
		mealPlanOptionVoteIDFetcher           func(*http.Request) string
		sessionContextDataFetcher             func(*http.Request) (*sessions.ContextData, error)
		useSearchService                      bool
	}
)

// ProvideService builds a new MealsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *textsearchcfg.Config,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	queueConfig *msgconfig.QueuesConfig,
	mealPlanningDataManager types.MealPlanningDataManager,
) (types.MealPlanningDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	searchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.MealSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeMeals)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing recipe index manager")
	}

	svc := &service{
		useSearchService:                      cfg.UseSearchService,
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:                     routeParamManager.BuildRouteParamStringIDFetcher(MealPlanIDURIParamKey),
		mealPlanGroceryListItemIDFetcher:      routeParamManager.BuildRouteParamStringIDFetcher(MealPlanGroceryListItemIDURIParamKey),
		mealPlanEventIDFetcher:                routeParamManager.BuildRouteParamStringIDFetcher(MealPlanEventIDURIParamKey),
		mealIDFetcher:                         routeParamManager.BuildRouteParamStringIDFetcher(MealIDURIParamKey),
		mealPlanOptionIDFetcher:               routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionIDURIParamKey),
		mealPlanTaskIDFetcher:                 routeParamManager.BuildRouteParamStringIDFetcher(MealPlanTaskIDURIParamKey),
		mealPlanOptionVoteIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionVoteIDURIParamKey),
		householdInstrumentOwnershipIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInstrumentOwnershipIDURIParamKey),
		userIngredientPreferenceIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(UserIngredientPreferenceIDURIParamKey),
		sessionContextDataFetcher:             sessions.FetchContextFromRequest,
		dataChangesPublisher:                  dataChangesPublisher,
		encoderDecoder:                        encoder,
		mealPlanningDataManager:               mealPlanningDataManager,
		tracer:                                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		mealsSearchIndex:                      searchIndex,
	}

	return svc, nil
}
