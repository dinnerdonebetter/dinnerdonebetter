package mealplantasks

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
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplans"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "meal_plan_tasks_service"
)

var _ types.MealPlanTaskDataService = (*service)(nil)

type (
	// service handles meal plan tasks.
	service struct {
		logger                    logging.Logger
		mealPlanTaskDataManager   types.MealPlanTaskDataManager
		mealPlanIDFetcher         func(*http.Request) string
		mealPlanTaskIDFetcher     func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new MealPlanTask.
func ProvideService(
	logger logging.Logger,
	mealPlanTaskDataManager types.MealPlanTaskDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.MealPlanTaskDataService, error) {
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
		mealPlanTaskIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(MealPlanTaskIDURIParamKey),
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		mealPlanTaskDataManager:   mealPlanTaskDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
