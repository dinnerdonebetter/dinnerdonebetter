package mealplanoptions

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "meal_plan_options_service"
)

var _ types.MealPlanOptionDataService = (*service)(nil)

type (
	// service handles meal plan options.
	service struct {
		logger                    logging.Logger
		mealPlanOptionDataManager types.MealPlanOptionDataManager
		mealPlanIDFetcher         func(*http.Request) string
		mealPlanEventIDFetcher    func(*http.Request) string
		mealPlanOptionIDFetcher   func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new MealPlanOptionsService.
func ProvideService(
	logger logging.Logger,
	mealPlanOptionDataManager types.MealPlanOptionDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.MealPlanOptionDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanEventIDFetcher:    routeParamManager.BuildRouteParamStringIDFetcher(mealplaneventsservice.MealPlanEventIDURIParamKey),
		mealPlanOptionIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionIDURIParamKey),
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		mealPlanOptionDataManager: mealPlanOptionDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
