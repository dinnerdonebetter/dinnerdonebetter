package mealplanfinalizer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningkeys "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers"
)

const (
	serviceName = "meal_plan_finalizer"
)

var _ workers.WorkerCounter = (*Worker)(nil)

type Worker struct {
	logger logging.Logger
	tracer tracing.Tracer

	dataManager             mealplanning.Repository
	postUpdatesPublisher    messagequeue.Publisher
	finalizedRecordsCounter metrics.Int64Counter
}

func NewMealPlanFinalizer(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	dataManager mealplanning.Repository,
	publisherProvider messagequeue.PublisherProvider,
	metricsProvider metrics.Provider,
	cfg *msgconfig.QueuesConfig,
) (*Worker, error) {
	finalizedRecordsCounter, err := metricsProvider.NewInt64Counter("meal_plan_finalizer.finalized_records")
	if err != nil {
		return nil, err
	}

	postUpdatesPublisher, err := publisherProvider.ProvidePublisher(ctx, cfg.DataChangesTopicName)
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
		changed, err = w.dataManager.AttemptToFinalizeMealPlan(ctx, mealPlan.ID, mealPlan.BelongsToAccount)
		if err != nil {
			return -1, observability.PrepareError(err, span, "finalizing meal plan")
		}

		if changed {
			changedCount++
			if err = w.postUpdatesPublisher.Publish(ctx, &audit.DataChangeMessage{
				Context: map[string]any{
					mealplanningkeys.MealPlanIDKey: mealPlan.ID,
					"meal_plan":                    mealPlan,
				},
				AccountID: mealPlan.BelongsToAccount,
			}); err != nil {
				logger.Error("writing data change message for finalized meal plan", err)
			}
		}
	}

	w.finalizedRecordsCounter.Add(ctx, changedCount)
	logger.WithValue("changed_count", changedCount).Info("finalized expired meal plans")

	return changedCount, nil
}
