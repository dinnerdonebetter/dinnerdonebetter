package grocerylistpreparation

import (
	"context"
	"fmt"

	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

// GroceryListCreator creates meal plan grocery lists for a given meal plan.
type GroceryListCreator interface {
	GenerateGroceryListInputs(ctx context.Context, mealPlan *types.MealPlan) ([]*types.MealPlanGroceryListItemDatabaseCreationInput, error)
}

type groceryListCreator struct {
	logger logging.Logger
	tracer tracing.Tracer
}

func NewGroceryListCreator(logger logging.Logger, tracerProvider tracing.TracerProvider) GroceryListCreator {
	return &groceryListCreator{
		logger: logging.EnsureLogger(logger).WithName("grocery_list_creator"),
		tracer: tracing.NewTracer(tracerProvider.Tracer("grocery_list_creator")),
	}
}

func (g *groceryListCreator) GenerateGroceryListInputs(ctx context.Context, mealPlan *types.MealPlan) ([]*types.MealPlanGroceryListItemDatabaseCreationInput, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	inputs := map[string]*types.MealPlanGroceryListItemDatabaseCreationInput{}
	logger := g.logger.Clone().WithValue(keys.MealPlanIDKey, mealPlan.ID)

	for _, event := range mealPlan.Events {
		logger = logger.WithValue(keys.MealPlanEventIDKey, event.ID)
		for _, option := range event.Options {
			if option.Chosen {
				logger = logger.WithValue(keys.MealPlanOptionIDKey, option.ID)
				for _, recipe := range option.Meal.Recipes {
					logger = logger.WithValue(keys.RecipeIDKey, recipe.ID)
					for _, step := range recipe.Steps {
						logger = logger.WithValue(keys.RecipeStepIDKey, step.ID)
						for _, ingredient := range step.Ingredients {
							if ingredient.Ingredient != nil {
								logger = logger.WithValue(keys.RecipeStepIngredientIDKey, ingredient.ID)
								if _, ok := inputs[ingredient.Ingredient.ID]; !ok {
									inputs[ingredient.Ingredient.ID] = &types.MealPlanGroceryListItemDatabaseCreationInput{
										Status:                 types.MealPlanGroceryListItemStatusUnknown,
										ValidMeasurementUnitID: ingredient.MeasurementUnit.ID,
										ValidIngredientID:      ingredient.Ingredient.ID,
										BelongsToMealPlan:      mealPlan.ID,
										ID:                     identifiers.New(),
										MinimumQuantityNeeded:  ingredient.MinimumQuantity,
										MaximumQuantityNeeded:  ingredient.MaximumQuantity,
									}
								} else {
									if inputs[ingredient.Ingredient.ID].ValidMeasurementUnitID == ingredient.MeasurementUnit.ID {
										inputs[ingredient.Ingredient.ID].MinimumQuantityNeeded += ingredient.MinimumQuantity
										inputs[ingredient.Ingredient.ID].MaximumQuantityNeeded += ingredient.MaximumQuantity
									} else {
										logger.Error(fmt.Errorf("mismatched measurement units: %s and %s", inputs[ingredient.Ingredient.ID].ValidMeasurementUnitID, ingredient.MeasurementUnit.ID), "creating grocery list")
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

	return dbInputs, nil
}
