package workers

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/recipeanalysis"
	"github.com/prixfixeco/api_server/pkg/types"
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
	customerDataCollector customerdata.Collector
}

// ProvideMealPlanGroceryListInitializer provides a MealPlanGroceryListInitializer.
func ProvideMealPlanGroceryListInitializer(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher recipeanalysis.RecipeAnalyzer,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *MealPlanGroceryListInitializer {
	return &MealPlanGroceryListInitializer{
		logger:                logging.EnsureLogger(logger).WithName(mealPlanGroceryListInitializerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(mealPlanGroceryListInitializerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		analyzer:              grapher,
		postUpdatesPublisher:  postUpdatesPublisher,
		customerDataCollector: customerDataCollector,
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

	var (
		errorResult                                  *multierror.Error
		mealPlansToGroceryListDatabaseCreationInputs = map[string][]*types.MealPlanGroceryListItemDatabaseCreationInput{}
	)

	for _, mealPlan := range mealPlans {
		inputs := map[string]*types.MealPlanGroceryListItemDatabaseCreationInput{}
		l := logger.Clone().WithValue(keys.MealPlanIDKey, mealPlan.ID)

		for _, event := range mealPlan.Events {
			l = l.WithValue(keys.MealPlanEventIDKey, event.ID)
			for _, option := range event.Options {
				if option.Chosen {
					l = l.WithValue(keys.MealPlanOptionIDKey, option.ID)
					for _, recipe := range option.Meal.Recipes {
						l = l.WithValue(keys.RecipeIDKey, recipe.ID)
						for _, step := range recipe.Steps {
							l = l.WithValue(keys.RecipeStepIDKey, step.ID)
							for _, ingredient := range step.Ingredients {
								if ingredient.Ingredient != nil {
									l = l.WithValue(keys.RecipeStepIngredientIDKey, ingredient.ID)
									if _, ok := inputs[ingredient.ID]; !ok {
										inputs[ingredient.ID] = &types.MealPlanGroceryListItemDatabaseCreationInput{
											Status:                 types.MealPlanGroceryListItemStatusUnknown,
											ValidMeasurementUnitID: ingredient.MeasurementUnit.ID,
											ValidIngredientID:      ingredient.Ingredient.ID,
											MealPlanOptionID:       option.ID,
											ID:                     ksuid.New().String(),
											MinimumQuantityNeeded:  ingredient.MinimumQuantity,
											MaximumQuantityNeeded:  ingredient.MaximumQuantity,
										}
									} else {
										if inputs[ingredient.ID].ValidMeasurementUnitID == ingredient.MeasurementUnit.ID {
											inputs[ingredient.ID].MinimumQuantityNeeded += ingredient.MinimumQuantity
											inputs[ingredient.ID].MaximumQuantityNeeded += ingredient.MaximumQuantity
										} else {
											l.Error(errors.New("mismatched measurement units"), "creating grocery list")
										}
									}
								}
							}
						}
					}
				}
			}
		}

		dbInputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{}
		for _, i := range inputs {
			dbInputs = append(dbInputs, i)
		}

		if err = w.dataManager.CreateMealPlanGroceryListItemsForMealPlan(ctx, mealPlan.ID, dbInputs); err != nil {
			errorResult = multierror.Append(errorResult, err)
		}

		mealPlansToGroceryListDatabaseCreationInputs[mealPlan.ID] = dbInputs
	}

	if errorResult == nil {
		return nil
	}

	return errorResult
}
