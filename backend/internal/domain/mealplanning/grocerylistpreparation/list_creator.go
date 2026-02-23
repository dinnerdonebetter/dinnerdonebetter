package grocerylistpreparation

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningkeys "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
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

// processRecipeIngredients processes ingredients from a recipe (main or associated) and adds them to the grocery list.
func (g *groceryListCreator) processRecipeIngredients(
	recipe *mealplanning.Recipe,
	recipeScale decimal.Decimal,
	optionID string,
	mealPlanID string,
	optionGroups map[string]int,
	selectionLookup map[string]uint16,
	aggregatedInputs map[string]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput,
	optionInputs *[]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput,
	logger logging.Logger,
) {
	for _, step := range recipe.Steps {
		logger = logger.WithValue(mealplanningkeys.RecipeStepIDKey, step.ID)
		for _, ingredient := range step.Ingredients {
			if ingredient.Ingredient == nil {
				continue
			}

			logger = logger.WithValue(mealplanningkeys.RecipeStepIngredientIDKey, ingredient.ID)

			minQty := float32(recipeScale.Mul(decimal.NewFromFloat32(ingredient.Quantity.Min)).Truncate(2).InexactFloat64())
			var maxQty *float32
			if ingredient.Quantity.Max != nil {
				maximum := float32(recipeScale.Mul(decimal.NewFromFloat32(*ingredient.Quantity.Max)).Truncate(2).InexactFloat64())
				maxQty = &maximum
			}

			optionGroupKey := fmt.Sprintf("%s:%d", step.ID, ingredient.Index)
			isOptionGroup := optionGroups[optionGroupKey] > 1

			if isOptionGroup {
				// This ingredient is part of an option group
				// Check if user made a selection for this option group
				selectionKey := fmt.Sprintf("%s:%d:%s", step.ID, ingredient.Index, mealplanning.MealPlanRecipeOptionSelectionTypeIngredient)
				selectedOptionIndex, hasSelection := selectionLookup[selectionKey]

				// Determine which option index to use:
				// - If user made a selection, use it
				// - If no selection, default to optionIndex=0 (the first/default option)
				targetOptionIndex := uint16(0)
				if hasSelection {
					targetOptionIndex = selectedOptionIndex
				}

				// Only add this ingredient if it matches the target option index
				if ingredient.OptionIndex != targetOptionIndex {
					continue
				}

				// Create grocery list item for the selected option with recipe context
				ingredientIndex := ingredient.Index
				optionIndex := ingredient.OptionIndex
				*optionInputs = append(*optionInputs, &mealplanning.MealPlanGroceryListItemDatabaseCreationInput{
					Status:                  mealplanning.MealPlanGroceryListItemStatusNeeds,
					ValidMeasurementUnitID:  ingredient.MeasurementUnit.ID,
					ValidIngredientID:       ingredient.Ingredient.ID,
					BelongsToMealPlan:       mealPlanID,
					BelongsToMealPlanOption: &optionID,
					RecipeID:                &recipe.ID,
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
						BelongsToMealPlanOption: &optionID,
						RecipeID:                &recipe.ID,
						RecipeStepID:            &step.ID,
						Status:                  mealplanning.MealPlanGroceryListItemStatusNeeds,
						ValidMeasurementUnitID:  ingredient.MeasurementUnit.ID,
						ValidIngredientID:       ingredient.Ingredient.ID,
						BelongsToMealPlan:       mealPlanID,
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

func (g *groceryListCreator) GenerateGroceryListInputs(ctx context.Context, mealPlan *mealplanning.MealPlan) ([]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	// Map to track option groups: key is (recipeStepID, ingredientIndex), value is count of options
	optionGroups := make(map[string]int)
	// Map for aggregated items (non-option items): key is ingredientID
	aggregatedInputs := make(map[string]*mealplanning.MealPlanGroceryListItemDatabaseCreationInput)
	// Slice for option items (items with alternatives): these are not aggregated
	optionInputs := []*mealplanning.MealPlanGroceryListItemDatabaseCreationInput{}

	logger := g.logger.Clone().WithValue(mealplanningkeys.MealPlanIDKey, mealPlan.ID)

	// Build a lookup map for user selections: key is (recipeStepID, ingredientIndex, selectionType), value is selectedOptionIndex
	selectionLookup := make(map[string]uint16)
	for _, selection := range mealPlan.Selections {
		if selection.SelectionType == mealplanning.MealPlanRecipeOptionSelectionTypeIngredient {
			selectionKey := fmt.Sprintf("%s:%d:%s", selection.RecipeStepID, selection.IngredientIndex, selection.SelectionType)
			selectionLookup[selectionKey] = selection.SelectedOptionIndex
		}
	}

	// First pass: identify option groups (ingredients with multiple options at the same index)
	// This includes both main recipes and their associated recipes
	for _, event := range mealPlan.Events {
		for _, option := range event.Options {
			if option.Chosen {
				for _, component := range option.Meal.Components {
					// Process main recipe
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
					// Process associated recipes
					for _, associatedRecipe := range component.Recipe.AssociatedRecipes {
						for _, step := range associatedRecipe.Steps {
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
	}

	// Second pass: process ingredients from main recipes and associated recipes
	for _, event := range mealPlan.Events {
		logger = logger.WithValue(mealplanningkeys.MealPlanEventIDKey, event.ID)
		for _, option := range event.Options {
			if option.Chosen {
				mealScale := decimal.NewFromFloat32(option.MealScale)
				logger = logger.WithValue(mealplanningkeys.MealPlanOptionIDKey, option.ID)
				for _, component := range option.Meal.Components {
					recipeScale := decimal.NewFromFloat32(component.RecipeScale).Mul(mealScale)
					logger = logger.WithValue(mealplanningkeys.RecipeIDKey, component.Recipe.ID)

					// Process main recipe ingredients
					g.processRecipeIngredients(
						&component.Recipe,
						recipeScale,
						option.ID,
						mealPlan.ID,
						optionGroups,
						selectionLookup,
						aggregatedInputs,
						&optionInputs,
						logger,
					)

					// Process associated recipe ingredients
					for _, associatedRecipe := range component.Recipe.AssociatedRecipes {
						logger = logger.WithValue(mealplanningkeys.RecipeIDKey, associatedRecipe.ID)
						g.processRecipeIngredients(
							associatedRecipe,
							recipeScale,
							option.ID,
							mealPlan.ID,
							optionGroups,
							selectionLookup,
							aggregatedInputs,
							&optionInputs,
							logger,
						)
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
