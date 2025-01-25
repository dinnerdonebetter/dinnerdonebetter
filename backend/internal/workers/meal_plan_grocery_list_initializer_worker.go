package workers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/businesslogic/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/hashicorp/go-multierror"
)

const (
	mealPlanGroceryListInitializerName = "meal_plan_grocery_list_initializer"
)

type (
	// MealPlanGroceryListInitializer initializes grocery lists for finalized meal plans.
	MealPlanGroceryListInitializer interface {
		InitializeGroceryListsForFinalizedMealPlans(ctx context.Context, _ []byte) error
	}

	// mealPlanGroceryListInitializer ensurers meal plan tasks are created.
	mealPlanGroceryListInitializer struct {
		logger               logging.Logger
		tracer               tracing.Tracer
		dataManager          database.DataManager
		postUpdatesPublisher messagequeue.Publisher
		groceryListCreator   grocerylistpreparation.GroceryListCreator
	}
)

// ProvideMealPlanGroceryListInitializer provides a mealPlanGroceryListInitializer.
func ProvideMealPlanGroceryListInitializer(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	tracerProvider tracing.TracerProvider,
	groceryListCreator grocerylistpreparation.GroceryListCreator,
) MealPlanGroceryListInitializer {
	return &mealPlanGroceryListInitializer{
		logger:               logging.EnsureLogger(logger).WithName(mealPlanGroceryListInitializerName),
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(mealPlanGroceryListInitializerName)),
		dataManager:          dataManager,
		postUpdatesPublisher: postUpdatesPublisher,
		groceryListCreator:   groceryListCreator,
	}
}

// InitializeGroceryListsForFinalizedMealPlans handles a pending write.
func (w *mealPlanGroceryListInitializer) InitializeGroceryListsForFinalizedMealPlans(ctx context.Context, _ []byte) error {
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

		l = l.WithValue("grocery_list_items_to_create", len(dbInputs))
		l.Info("creating grocery list items for meal plan")

		for _, dbInput := range dbInputs {
			var createdItem *types.MealPlanGroceryListItem
			createdItem, err = w.dataManager.CreateMealPlanGroceryListItem(ctx, dbInput)
			if err != nil {
				errorResult = multierror.Append(errorResult, err)
				l.Error("failed to create grocery list for meal plan", err)
				continue
			}

			if err = w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				MealPlanGroceryListItem:   createdItem,
				MealPlanGroceryListItemID: createdItem.ID,
				EventType:                 types.MealPlanGroceryListItemCreatedServiceEventType,
				MealPlanID:                dbInput.BelongsToMealPlan,
			}); err != nil {
				l.Error("failed to write update message for meal plan grocery list item", err)
			}
		}
	}

	return errorResult.ErrorOrNil()
}
