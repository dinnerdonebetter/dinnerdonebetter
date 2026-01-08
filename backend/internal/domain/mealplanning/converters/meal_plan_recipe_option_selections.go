package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

// ConvertMealPlanRecipeOptionSelectionToMealPlanRecipeOptionSelectionDatabaseCreationInput builds a MealPlanRecipeOptionSelectionDatabaseCreationInput from a MealPlanRecipeOptionSelection.
func ConvertMealPlanRecipeOptionSelectionToMealPlanRecipeOptionSelectionDatabaseCreationInput(input *mealplanning.MealPlanRecipeOptionSelection) *mealplanning.MealPlanRecipeOptionSelectionDatabaseCreationInput {
	return &mealplanning.MealPlanRecipeOptionSelectionDatabaseCreationInput{
		ID:                      input.ID,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		RecipeID:                input.RecipeID,
		RecipeStepID:            input.RecipeStepID,
		IngredientIndex:         input.IngredientIndex,
		SelectedOptionIndex:     input.SelectedOptionIndex,
		SelectionType:           input.SelectionType,
	}
}

// ConvertMealPlanRecipeOptionSelectionDatabaseCreationInputToMealPlanRecipeOptionSelectionDatabaseCreationInput creates a new DatabaseCreationInput with a new ID.
func ConvertMealPlanRecipeOptionSelectionDatabaseCreationInputToMealPlanRecipeOptionSelectionDatabaseCreationInput(input *mealplanning.MealPlanRecipeOptionSelectionDatabaseCreationInput) *mealplanning.MealPlanRecipeOptionSelectionDatabaseCreationInput {
	return &mealplanning.MealPlanRecipeOptionSelectionDatabaseCreationInput{
		ID:                      identifiers.New(),
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		RecipeID:                input.RecipeID,
		RecipeStepID:            input.RecipeStepID,
		IngredientIndex:         input.IngredientIndex,
		SelectedOptionIndex:     input.SelectedOptionIndex,
		SelectionType:           input.SelectionType,
	}
}
