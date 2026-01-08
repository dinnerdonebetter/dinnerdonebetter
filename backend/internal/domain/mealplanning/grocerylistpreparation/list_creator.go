package grocerylistpreparation

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	"github.com/shopspring/decimal"
)

// GroceryListCreator creates meal plan grocery lists for a given meal plan.
type GroceryListCreator interface {
	GenerateGroceryListInputs(ctx context.Context, mealPlan *mealplanning.MealPlan) ([]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput, error)
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

func (g *groceryListCreator) GenerateGroceryListInputs(ctx context.Context, mealPlan *mealplanning.MealPlan) ([]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	// Map to track option groups: key is (recipeStepID, ingredientIndex), value is count of options
	optionGroups := make(map[string]int)
	// Map for aggregated items (non-option items): key is ingredientID
	aggregatedInputs := make(map[string]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput)
	// Slice for option items (items with alternatives): these are not aggregated
	optionInputs := []*mealplanning.MealPlanGroceryListItemDatabaseCreationInput{}

	logger := g.logger.Clone().WithValue(keys.MealPlanIDKey, mealPlan.ID)

	// First pass: identify option groups (ingredients with multiple options at the same index)
	for _, event := range mealPlan.Events {
		for _, option := range event.Options {
			if option.Chosen {
				for _, component := range option.Meal.Components {
					for _, step := range component.Recipe.Steps {
						// Track how many options exist for each (stepID, index) combination
						indexCounts := make(map[uint16]int)
						for _, ingredient := range step.Ingredients {
							if ingredient.Ingredient != nil {
								indexCounts[ingredient.Index]++
							}
						}
						// Mark groups with multiple options
						for index, count := range indexCounts {
							if count > 1 {
								optionGroupKey := fmt.Sprintf("%s:%d", step.ID, index)
								optionGroups[optionGroupKey] = count
							}
						}
					}
				}
			}
		}
	}

	// Second pass: process ingredients
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
							if ingredient.Ingredient == nil {
								continue
							}

							logger = logger.WithValue(keys.RecipeStepIngredientIDKey, ingredient.ID)

							minQty := float32(recipeScale.Mul(decimal.NewFromFloat32(ingredient.Quantity.Min)).Truncate(2).InexactFloat64())
							var maxQty *float32
							if ingredient.Quantity.Max != nil {
								maximum := float32(recipeScale.Mul(decimal.NewFromFloat32(*ingredient.Quantity.Max)).Truncate(2).InexactFloat64())
								maxQty = &maximum
							}

							optionGroupKey := fmt.Sprintf("%s:%d", step.ID, ingredient.Index)
							isOptionGroup := optionGroups[optionGroupKey] > 1

							if isOptionGroup {
								// This ingredient is part of an option group - create separate item with recipe context
								ingredientIndex := ingredient.Index
								optionIndex := ingredient.OptionIndex
								optionInputs = append(optionInputs, &mealplanning.MealPlanGroceryListItemDatabaseCreationInput{
									Status:                  mealplanning.MealPlanGroceryListItemStatusNeeds,
									ValidMeasurementUnitID:  ingredient.MeasurementUnit.ID,
									ValidIngredientID:       ingredient.Ingredient.ID,
									BelongsToMealPlan:       mealPlan.ID,
									BelongsToMealPlanOption: &option.ID,
									RecipeID:                &component.Recipe.ID,
									RecipeStepID:            &step.ID,
									IngredientIndex:         &ingredientIndex,
									OptionIndex:             &optionIndex,
									ID:                      identifiers.New(),
									QuantityNeeded: types.Float32RangeWithOptionalMax{
										Max: maxQty,
										Min: minQty,
									},
								})
							} else {
								// This ingredient is not part of an option group - aggregate as before
								if existing, ok := aggregatedInputs[ingredient.Ingredient.ID]; !ok {
									aggregatedInputs[ingredient.Ingredient.ID] = &mealplanning.MealPlanGroceryListItemDatabaseCreationInput{
										BelongsToMealPlanOption: &option.ID,
										RecipeID:                &component.Recipe.ID,
										RecipeStepID:            &step.ID,
										Status:                  mealplanning.MealPlanGroceryListItemStatusNeeds,
										ValidMeasurementUnitID:  ingredient.MeasurementUnit.ID,
										ValidIngredientID:       ingredient.Ingredient.ID,
										BelongsToMealPlan:       mealPlan.ID,
										ID:                      identifiers.New(),
										QuantityNeeded: types.Float32RangeWithOptionalMax{
											Max: maxQty,
											Min: minQty,
										},
									}
								} else {
									if existing.ValidMeasurementUnitID == ingredient.MeasurementUnit.ID {
										existing.QuantityNeeded.Min += minQty

										if existing.QuantityNeeded.Max != nil && maxQty != nil {
											*existing.QuantityNeeded.Max += *maxQty
										} else if maxQty != nil {
											existing.QuantityNeeded.Max = maxQty
										}
									} else {
										logger.Error("creating grocery list", fmt.Errorf("mismatched measurement units: %s and %s", existing.ValidMeasurementUnitID, ingredient.MeasurementUnit.ID))
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Combine aggregated items and option items
	dbInputs := make([]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput, 0, len(aggregatedInputs)+len(optionInputs))
	for _, i := range aggregatedInputs {
		dbInputs = append(dbInputs, i)
	}
	dbInputs = append(dbInputs, optionInputs...)

	return dbInputs, nil
}
