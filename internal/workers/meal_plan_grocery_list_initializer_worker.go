package workers

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/features/grocerylistpreparation"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

const (
	mealPlanGroceryListInitializerName = "meal_plan_grocery_list_initializer"
)

// MealPlanGroceryListInitializer ensurers meal plan tasks are created.
type MealPlanGroceryListInitializer struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	analyzer              recipeanalysis.RecipeAnalyzer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	customerDataCollector analytics.EventReporter
	groceryListCreator    grocerylistpreparation.GroceryListCreator
}

// ProvideMealPlanGroceryListInitializer provides a MealPlanGroceryListInitializer.
func ProvideMealPlanGroceryListInitializer(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher recipeanalysis.RecipeAnalyzer,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector analytics.EventReporter,
	tracerProvider tracing.TracerProvider,
	groceryListCreator grocerylistpreparation.GroceryListCreator,
) *MealPlanGroceryListInitializer {
	return &MealPlanGroceryListInitializer{
		logger:                logging.EnsureLogger(logger).WithName(mealPlanGroceryListInitializerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(mealPlanGroceryListInitializerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		analyzer:              grapher,
		postUpdatesPublisher:  postUpdatesPublisher,
		customerDataCollector: customerDataCollector,
		groceryListCreator:    groceryListCreator,
	}
}

// HandleMessage handles a pending write.
func (w *MealPlanGroceryListInitializer) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlans, err := w.dataManager.GetFinalizedMealPlansWithUninitializedGroceryLists(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting finalized meal plan data")
	}

	var errorResult *multierror.Error

	for _, mealPlan := range mealPlans {
		dbInputs, groceryListCreationErr := w.groceryListCreator.GenerateGroceryListInputs(ctx, mealPlan)
		if groceryListCreationErr != nil {
			errorResult = multierror.Append(errorResult, groceryListCreationErr)
		}

		if err = w.dataManager.CreateMealPlanGroceryListItemsForMealPlan(ctx, mealPlan.ID, dbInputs); err != nil {
			errorResult = multierror.Append(errorResult, err)
		}
	}

	if errorResult == nil {
		return nil
	}

	return errorResult
}
