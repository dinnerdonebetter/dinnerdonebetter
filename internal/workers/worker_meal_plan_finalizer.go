package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	name = "meal_plan_finalizer"
)

// MealPlanFinalizationWorker finalizes meal plans.
type MealPlanFinalizationWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	emailSender           email.Emailer
	customerDataCollector customerdata.Collector
}

// ProvideMealPlanFinalizationWorker provides a MealPlanFinalizationWorker.
func ProvideMealPlanFinalizationWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *MealPlanFinalizationWorker {
	return &MealPlanFinalizationWorker{
		logger:                logging.EnsureLogger(logger).WithName(name),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		postUpdatesPublisher:  postUpdatesPublisher,
		emailSender:           emailSender,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a message ordering the finalization of expired meal plans.
func (w *MealPlanFinalizationWorker) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	return w.finalizeExpiredMealPlans(ctx)
}

func (w *MealPlanFinalizationWorker) finalizeExpiredMealPlans(ctx context.Context) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlans, fetchMealPlansErr := w.dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if fetchMealPlansErr != nil {
		return observability.PrepareAndLogError(fetchMealPlansErr, logger, span, "fetching unfinalized and expired meal plan")
	}

	logger.WithValue("quantity", len(mealPlans)).Info("finalizing expired meal plans")

	for _, mealPlan := range mealPlans {
		changed, err := w.dataManager.AttemptToFinalizeMealPlan(ctx, mealPlan.ID, mealPlan.BelongsToHousehold)
		if err != nil {
			return observability.PrepareError(err, span, "finalizing meal plan")
		}

		if !changed {
			logger.Debug("meal plan was not changed")
		}
	}

	return nil
}
