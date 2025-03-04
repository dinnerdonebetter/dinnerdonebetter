package mealplangrocerylistinitializer

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/workers"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/hashicorp/go-multierror"
)

const (
	serviceName = "meal_plan_grocery_list_initializer"
)

var _ workers.Worker = (*Worker)(nil)

type Worker struct {
	logger                  logging.Logger
	tracer                  tracing.Tracer
	dataManager             database.DataManager // TODO: make this less potent
	postUpdatesPublisher    messagequeue.Publisher
	recordsProcessedCounter metrics.Int64Counter
	// TODO: groceryListCreator      grocerylistpreparation.GroceryListCreator
}

func NewMealPlanGroceryListInitializer(logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	publisherProvider messagequeue.PublisherProvider,
	// TODO: groceryListCreator grocerylistpreparation.GroceryListCreator,
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
		// TODO: groceryListCreator:      groceryListCreator,
		logger: logging.EnsureLogger(logger).WithName(serviceName),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
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

		dbInputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{}
		/*
			// TODO: RESTOREME
			//	dbInputs, err = w.groceryListCreator.GenerateGroceryListInputs(ctx, mealPlan)
			//	if err != nil {
			//		errorResult = multierror.Append(errorResult, err)
			//		l.Error("failed to generate grocery list inputs for meal plan", err)
			//		continue
			//	}
		*/

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

			if err = w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				MealPlanGroceryListItem:   createdItem,
				MealPlanGroceryListItemID: createdItem.ID,
				EventType:                 types.MealPlanGroceryListItemCreatedServiceEventType,
				MealPlanID:                dbInput.BelongsToMealPlan,
			}); err != nil {
				l.Error("failed to write update message for meal plan grocery list item", err)
			}
		}
		w.recordsProcessedCounter.Add(ctx, createdCount)
	}

	return errorResult.ErrorOrNil()
}
