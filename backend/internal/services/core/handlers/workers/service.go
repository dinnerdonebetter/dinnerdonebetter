package workers

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/internal/services/eating/businesslogic/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/services/eating/workers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "workers_service"
)

var _ types.WorkerService = (*service)(nil)

type (
	// service handles worker invocation requests.
	service struct {
		logger                         logging.Logger
		sessionContextDataFetcher      func(*http.Request) (*sessioncontext.SessionContextData, error)
		encoderDecoder                 encoding.ServerEncoderDecoder
		tracer                         tracing.Tracer
		mealPlanFinalizationWorker     workers.MealPlanFinalizationWorker
		mealPlanGroceryListInitializer workers.MealPlanGroceryListInitializer
		mealPlanTaskCreatorWorker      workers.MealPlanTaskCreatorWorker
	}
)

// ProvideService builds a new WorkerService.
func ProvideService(
	logger logging.Logger,
	dataManager database.DataManager,
	encoder encoding.ServerEncoderDecoder,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	grapher recipeanalysis.RecipeAnalyzer,
	queueConfig *msgconfig.QueuesConfig,
) (types.WorkerService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	mealPlanFinalizationWorker, err := workers.ProvideMealPlanFinalizationWorker(
		logger,
		dataManager,
		dataChangesPublisher,
		tracerProvider,
		metricsProvider,
	)
	if err != nil {
		return nil, fmt.Errorf("setting up %s meal plan finalization worker: %w", serviceName, err)
	}

	mealPlanGroceryListInitializer := workers.ProvideMealPlanGroceryListInitializer(
		logger,
		dataManager,
		dataChangesPublisher,
		tracerProvider,
		grocerylistpreparation.NewGroceryListCreator(logger, tracerProvider),
	)

	mealPlanTaskCreatorWorker := workers.ProvideMealPlanTaskCreationEnsurerWorker(
		logger,
		dataManager,
		grapher,
		dataChangesPublisher,
		tracerProvider,
	)

	svc := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:                 encoder,
		mealPlanFinalizationWorker:     mealPlanFinalizationWorker,
		mealPlanGroceryListInitializer: mealPlanGroceryListInitializer,
		mealPlanTaskCreatorWorker:      mealPlanTaskCreatorWorker,
		sessionContextDataFetcher:      sessioncontext.FetchContextFromRequest,
		tracer:                         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
