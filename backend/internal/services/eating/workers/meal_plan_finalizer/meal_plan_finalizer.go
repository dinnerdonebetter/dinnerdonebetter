package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/workers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName = "meal_plan_finalizer"
)

var _ workers.WorkerCounter = (*Worker)(nil)

type Worker struct {
	logger logging.Logger
	tracer tracing.Tracer

	dataManager             database.DataManager
	postUpdatesPublisher    messagequeue.Publisher
	finalizedRecordsCounter metrics.Int64Counter
}

func NewMealPlanFinalizer(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	dataManager database.DataManager,
	publisherProvider messagequeue.PublisherProvider,
	metricsProvider metrics.Provider,
	cfg *msgconfig.QueuesConfig,
) (*Worker, error) {
	finalizedRecordsCounter, err := metricsProvider.NewInt64Counter("meal_plan_finalizer.finalized_records")
	if err != nil {
		return nil, err
	}

	postUpdatesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, err
	}

	return &Worker{
		dataManager:             dataManager,
		postUpdatesPublisher:    postUpdatesPublisher,
		finalizedRecordsCounter: finalizedRecordsCounter,

		logger: logging.EnsureLogger(logger).WithName(serviceName),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}, nil
}

func (w *Worker) Work(ctx context.Context) (int64, error) {
	ctx, span := w.tracer.StartSpan(ctx)
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

	var changedCount int64
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

	w.finalizedRecordsCounter.Add(ctx, changedCount)
	logger.WithValue("changed_count", changedCount).Info("finalized expired meal plans")

	return changedCount, nil
}
