package workers

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName = "meal_plan_finalizer"
)

type (
	MealPlanFinalizationWorker interface {
		FinalizeExpiredMealPlans(ctx context.Context) (int, error)
		FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context) error
	}

	// mealPlanFinalizationWorker finalizes meal plans.
	mealPlanFinalizationWorker struct {
		logger                  logging.Logger
		tracer                  tracing.Tracer
		dataManager             database.DataManager
		postUpdatesPublisher    messagequeue.Publisher
		finalizedRecordsCounter metrics.Int64Counter
	}
)

// ProvideMealPlanFinalizationWorker provides a mealPlanFinalizationWorker.
func ProvideMealPlanFinalizationWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
) (MealPlanFinalizationWorker, error) {
	counter, err := metricsProvider.NewInt64Counter(fmt.Sprintf("%s.finalized_records", serviceName))
	if err != nil {
		return nil, fmt.Errorf("failed to create counter for finalized records metrics: %w", err)
	}

	return &mealPlanFinalizationWorker{
		logger:                  logging.EnsureLogger(logger).WithName(serviceName),
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataManager:             dataManager,
		postUpdatesPublisher:    postUpdatesPublisher,
		finalizedRecordsCounter: counter,
	}, nil
}

// finalizeExpiredMealPlans handles a message ordering the finalization of expired meal plans.
func (w *mealPlanFinalizationWorker) finalizeExpiredMealPlans(ctx context.Context) (int, error) {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()
	logger.Info("beginning finalization of expired meal plans")

	mealPlans, err := w.dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if err != nil {
		return -1, observability.PrepareAndLogError(err, logger, span, "fetching unfinalized and expired meal plan")
	}

	if len(mealPlans) > 0 {
		logger.WithValue("quantity", len(mealPlans)).Info("finalizing expired meal plans")
	}

	var changedCount int
	for _, mealPlan := range mealPlans {
		var changed bool
		changed, err = w.dataManager.AttemptToFinalizeMealPlan(ctx, mealPlan.ID, mealPlan.BelongsToHousehold)
		if err != nil {
			return -1, observability.PrepareError(err, span, "finalizing meal plan")
		}

		if changed {
			changedCount++
			if err = w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				MealPlanID:  mealPlan.ID,
				MealPlan:    mealPlan,
				HouseholdID: mealPlan.BelongsToHousehold,
			}); err != nil {
				logger.Error("writing data change message for finalized meal plan", err)
			}
		}
	}

	w.finalizedRecordsCounter.Add(ctx, int64(changedCount))
	logger.WithValue("changed_count", changedCount).Info("finalized expired meal plans")

	return changedCount, nil
}

// FinalizeExpiredMealPlans handles a message ordering the finalization of expired meal plans.
func (w *mealPlanFinalizationWorker) FinalizeExpiredMealPlans(ctx context.Context) (int, error) {
	return w.finalizeExpiredMealPlans(ctx)
}

// FinalizeExpiredMealPlansWithoutReturningCount handles a message ordering the finalization of expired meal plans.
func (w *mealPlanFinalizationWorker) FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context) error {
	_, err := w.finalizeExpiredMealPlans(ctx)
	return err
}
