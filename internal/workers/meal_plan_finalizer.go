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
	"github.com/prixfixeco/api_server/pkg/types"
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
	n := "meal_plan_finalizer"

	return &MealPlanFinalizationWorker{
		logger:                logging.EnsureLogger(logger).WithName(n),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(n)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		postUpdatesPublisher:  postUpdatesPublisher,
		emailSender:           emailSender,
		customerDataCollector: customerDataCollector,
	}
}

// finalizeExpiredMealPlans handles a message ordering the finalization of expired meal plans.
func (w *MealPlanFinalizationWorker) finalizeExpiredMealPlans(ctx context.Context) (int, error) {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlans, fetchMealPlansErr := w.dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if fetchMealPlansErr != nil {
		return -1, observability.PrepareAndLogError(fetchMealPlansErr, logger, span, "fetching unfinalized and expired meal plan")
	}

	logger.WithValue("quantity", len(mealPlans)).Info("finalizing expired meal plans")

	var changedCount int
	for _, mealPlan := range mealPlans {
		changed, err := w.dataManager.AttemptToFinalizeMealPlan(ctx, mealPlan.ID, mealPlan.BelongsToHousehold)
		if err != nil {
			return -1, observability.PrepareError(err, span, "finalizing meal plan")
		}

		if changed {
			changedCount++

			if dataChangePublishErr := w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				DataType:                  types.MealPlanDataType,
				EventType:                 types.MealPlanFinalizedCustomerEventType,
				MealPlan:                  mealPlan,
				MealPlanID:                mealPlan.ID,
				Context:                   map[string]string{},
				AttributableToHouseholdID: mealPlan.BelongsToHousehold,
			}); dataChangePublishErr != nil {
				observability.AcknowledgeError(dataChangePublishErr, logger, span, "publishing data change message")
			}
		}
	}

	return changedCount, nil
}

// FinalizeExpiredMealPlans handles a message ordering the finalization of expired meal plans.
func (w *MealPlanFinalizationWorker) FinalizeExpiredMealPlans(ctx context.Context, _ []byte) (int, error) {
	return w.finalizeExpiredMealPlans(ctx)
}

// FinalizeExpiredMealPlansWithoutReturningCount handles a message ordering the finalization of expired meal plans.
func (w *MealPlanFinalizationWorker) FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context, _ []byte) error {
	_, err := w.finalizeExpiredMealPlans(ctx)
	return err
}
