package grocerylistpreparation

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/shopspring/decimal"
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
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("grocery_list_creator")),
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
				mealScale := decimal.NewFromFloat32(option.MealScale)
				logger = logger.WithValue(keys.MealPlanOptionIDKey, option.ID)
				for _, component := range option.Meal.Components {
					recipeScale := decimal.NewFromFloat32(component.RecipeScale).Mul(mealScale)
					logger = logger.WithValue(keys.RecipeIDKey, component.Recipe.ID)
					for _, step := range component.Recipe.Steps {
						logger = logger.WithValue(keys.RecipeStepIDKey, step.ID)
						for _, ingredient := range step.Ingredients {
							if ingredient.Ingredient != nil {
								logger = logger.WithValue(keys.RecipeStepIngredientIDKey, ingredient.ID)
								if _, ok := inputs[ingredient.Ingredient.ID]; !ok {
									minQty := float32(recipeScale.Mul(decimal.NewFromFloat32(ingredient.Quantity.Min)).Truncate(2).InexactFloat64())
									var maxQty *float32
									if ingredient.Quantity.Max != nil {
										maximum := float32(recipeScale.Mul(decimal.NewFromFloat32(*ingredient.Quantity.Max)).Truncate(2).InexactFloat64())
										maxQty = &maximum
									}

									inputs[ingredient.Ingredient.ID] = &types.MealPlanGroceryListItemDatabaseCreationInput{
										Status:                 types.MealPlanGroceryListItemStatusNeeds,
										ValidMeasurementUnitID: ingredient.MeasurementUnit.ID,
										ValidIngredientID:      ingredient.Ingredient.ID,
										BelongsToMealPlan:      mealPlan.ID,
										ID:                     identifiers.New(),
										QuantityNeeded: types.Float32RangeWithOptionalMax{
											Max: maxQty,
											Min: minQty,
										},
									}
								} else {
									if inputs[ingredient.Ingredient.ID].ValidMeasurementUnitID == ingredient.MeasurementUnit.ID {
										inputs[ingredient.Ingredient.ID].QuantityNeeded.Min += ingredient.Quantity.Min

										if inputs[ingredient.Ingredient.ID].QuantityNeeded.Max != nil {
											if ingredient.Quantity.Max != nil {
												*inputs[ingredient.Ingredient.ID].QuantityNeeded.Max += *ingredient.Quantity.Max
											}
										} else if ingredient.Quantity.Max != nil {
											inputs[ingredient.Ingredient.ID].QuantityNeeded.Max = ingredient.Quantity.Max
										}
									} else {
										logger.Error("creating grocery list", fmt.Errorf("mismatched measurement units: %s and %s", inputs[ingredient.Ingredient.ID].ValidMeasurementUnitID, ingredient.MeasurementUnit.ID))
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
