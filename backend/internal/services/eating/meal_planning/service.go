package mealplanning

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	textsearch "github.com/dinnerdonebetter/backend/internal/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
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
		searchIndex                           textsearch.IndexSearcher[types.MealSearchSubset]
		householdInstrumentOwnershipIDFetcher func(*http.Request) string
		userIngredientPreferenceIDFetcher     func(*http.Request) string
		mealPlanIDFetcher                     func(*http.Request) string
		mealPlanGroceryListItemIDFetcher      func(*http.Request) string
		mealPlanEventIDFetcher                func(*http.Request) string
		mealIDFetcher                         func(*http.Request) string
		mealPlanOptionIDFetcher               func(*http.Request) string
		mealPlanTaskIDFetcher                 func(*http.Request) string
		mealPlanOptionVoteIDFetcher           func(*http.Request) string
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
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

	searchIndex, err := textsearchcfg.ProvideIndex[types.MealSearchSubset](ctx, logger, tracerProvider, searchConfig, textsearch.IndexTypeMeals)
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
		sessionContextDataFetcher:             authentication.FetchContextFromRequest,
		dataChangesPublisher:                  dataChangesPublisher,
		encoderDecoder:                        encoder,
		mealPlanningDataManager:               mealPlanningDataManager,
		tracer:                                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:                           searchIndex,
	}

	return svc, nil
}
