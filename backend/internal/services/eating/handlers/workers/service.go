package workers

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/workers"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/eating/workers/meal_plan_task_creator"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "workers_service"
)

var _ types.WorkerService = (*service)(nil)

type (
	// service handles worker invocation requests.
	service struct {
		logger                               logging.Logger
		sessionContextDataFetcher            func(*http.Request) (*sessions.ContextData, error)
		encoderDecoder                       encoding.ServerEncoderDecoder
		tracer                               tracing.Tracer
		mealPlanFinalizerWorker              workers.WorkerCounter
		mealPlanGroceryListInitializerWorker workers.Worker
		mealPlanTaskCreatorWorker            workers.Worker
	}
)

// ProvideService builds a new WorkerService.
func ProvideService(
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	mealPlanFinalizerWorker *mealplanfinalizer.Worker,
	mealPlanGroceryListInitializerWorker *mealplangrocerylistinitializer.Worker,
	mealPlanTaskCreatorWorker *mealplantaskcreator.Worker,
) (types.WorkerService, error) {
	svc := &service{
		logger:                               logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:                       encoder,
		mealPlanFinalizerWorker:              mealPlanFinalizerWorker,
		mealPlanGroceryListInitializerWorker: mealPlanGroceryListInitializerWorker,
		mealPlanTaskCreatorWorker:            mealPlanTaskCreatorWorker,
		sessionContextDataFetcher:            sessions.FetchContextFromRequest,
		tracer:                               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
