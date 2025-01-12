package mealplanoptionvotes

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanevents"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptions"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplans"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "meal_plan_option_votes_service"
)

var _ types.MealPlanOptionVoteDataService = (*service)(nil)

type (
	dataManager interface {
		types.MealPlanDataManager
		types.MealPlanEventDataManager
		types.MealPlanOptionDataManager
		types.MealPlanOptionVoteDataManager
	}

	// service handles meal plan option votes.
	service struct {
		logger                      logging.Logger
		dataManager                 dataManager
		mealPlanIDFetcher           func(*http.Request) string
		mealPlanEventIDFetcher      func(*http.Request) string
		mealPlanOptionIDFetcher     func(*http.Request) string
		mealPlanOptionVoteIDFetcher func(*http.Request) string
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher        messagequeue.Publisher
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
	}
)

// ProvideService builds a new MealPlanOptionVotesService.
func ProvideService(
	logger logging.Logger,
	dataManager database.DataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.MealPlanOptionVoteDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanEventIDFetcher:      routeParamManager.BuildRouteParamStringIDFetcher(mealplaneventsservice.MealPlanEventIDURIParamKey),
		mealPlanOptionIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(mealplanoptionsservice.MealPlanOptionIDURIParamKey),
		mealPlanOptionVoteIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionVoteIDURIParamKey),
		sessionContextDataFetcher:   authentication.FetchContextFromRequest,
		dataManager:                 dataManager,
		dataChangesPublisher:        dataChangesPublisher,
		encoderDecoder:              encoder,
		tracer:                      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
