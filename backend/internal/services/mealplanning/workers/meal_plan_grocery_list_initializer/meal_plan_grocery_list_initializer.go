package mealplangrocerylistinitializer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/mealplanning/businesslogic/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers"

	"github.com/hashicorp/go-multierror"
)

const (
	serviceName = "meal_plan_grocery_list_initializer"
)

var _ workers.Worker = (*Worker)(nil)

type Worker struct {
	logger                  logging.Logger
	tracer                  tracing.Tracer
	dataManager             types.Repository // TODO: make this less potent
	postUpdatesPublisher    messagequeue.Publisher
	recordsProcessedCounter metrics.Int64Counter
	groceryListCreator      grocerylistpreparation.GroceryListCreator
}

func NewMealPlanGroceryListInitializer(logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	publisherProvider messagequeue.PublisherProvider,
	groceryListCreator grocerylistpreparation.GroceryListCreator,
	cfg *msgconfig.QueuesConfig,
) (*Worker, error) {
	postUpdatesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, err
	}

	recordsProcessedCounter, err := metricsProvider.NewInt64Counter("meal_plan_grocery_list_initializer.records_processed")
	if err != nil {
		return nil, err
	}

	return &Worker{
		recordsProcessedCounter: recordsProcessedCounter,
		postUpdatesPublisher:    postUpdatesPublisher,
		groceryListCreator:      groceryListCreator,
		logger:                  logging.EnsureLogger(logger).WithName(serviceName),
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}, nil
}

func (w *Worker) Work(ctx context.Context) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlans, err := w.dataManager.GetFinalizedMealPlansWithUninitializedGroceryLists(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting finalized meal plan data")
	}

	logger = logger.WithValue("meal_plan_quantity", len(mealPlans))

	if len(mealPlans) > 0 {
		logger.Info("attempting to initialize grocery lists for meal plans")
	}

	errorResult := &multierror.Error{}

	for _, mealPlan := range mealPlans {
		l := logger.WithValue(keys.MealPlanIDKey, mealPlan.ID)

		var dbInputs []*types.MealPlanGroceryListItemDatabaseCreationInput
		dbInputs, err = w.groceryListCreator.GenerateGroceryListInputs(ctx, mealPlan)
		if err != nil {
			errorResult = multierror.Append(errorResult, err)
			l.Error("failed to generate grocery list inputs for meal plan", err)
			continue
		}

		l = l.WithValue("to_create", len(dbInputs))
		l.Info("creating grocery list items for meal plan")

		var createdCount int64
		for _, dbInput := range dbInputs {
			var createdItem *types.MealPlanGroceryListItem
			createdItem, err = w.dataManager.CreateMealPlanGroceryListItem(ctx, dbInput)
			if err != nil {
				errorResult = multierror.Append(errorResult, err)
				l.Error("failed to create grocery list for meal plan", err)
				continue
			}
			createdCount++

			if err = w.postUpdatesPublisher.Publish(ctx, &audit.DataChangeMessage{
				EventType: types.MealPlanGroceryListItemCreatedServiceEventType,
				Context: map[string]any{
					"groceryListItem":                 createdItem,
					keys.MealPlanGroceryListItemIDKey: createdItem.ID,
					keys.MealPlanIDKey:                dbInput.BelongsToMealPlan,
				},
			}); err != nil {
				l.Error("failed to write update message for meal plan grocery list item", err)
			}
		}
		w.recordsProcessedCounter.Add(ctx, createdCount)
	}

	return errorResult.ErrorOrNil()
}
