package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionConditionDatabaseCreationInput creates a RecipeStepCompletionConditionDatabaseCreationInput from a RecipeStepCompletionConditionCreationRequestInput.
func ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionConditionDatabaseCreationInput(input *types.RecipeStepCompletionConditionCreationRequestInput, recipeStep *types.RecipeStepDatabaseCreationInput) *types.RecipeStepCompletionConditionDatabaseCreationInput {
	id := identifiers.New()

	var ingredients []*types.RecipeStepCompletionConditionIngredientDatabaseCreationInput
	for _, i := range input.Ingredients {
		x := ConvertRecipeStepCompletionConditionIngredientCreationRequestInputToRecipeStepCompletionConditionIngredientDatabaseCreationInput(i, recipeStep)
		x.BelongsToRecipeStepCompletionCondition = id
		ingredients = append(ingredients, x)
	}

	y := &types.RecipeStepCompletionConditionDatabaseCreationInput{
		ID:                  id,
		IngredientStateID:   input.IngredientStateID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Ingredients:         ingredients,
		Optional:            input.Optional,
	}

	return y
}

// ConvertRecipeStepCompletionConditionIngredientCreationRequestInputToRecipeStepCompletionConditionIngredientDatabaseCreationInput creates a RecipeStepCompletionConditionIngredientDatabaseCreationInput from a RecipeStepCompletionConditionCreationRequestInput.
func ConvertRecipeStepCompletionConditionIngredientCreationRequestInputToRecipeStepCompletionConditionIngredientDatabaseCreationInput(input *types.RecipeStepCompletionConditionIngredientCreationRequestInput, recipeStep *types.RecipeStepDatabaseCreationInput) *types.RecipeStepCompletionConditionIngredientDatabaseCreationInput {
	x := &types.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
		ID:                   identifiers.New(),
		RecipeStepIngredient: recipeStep.Ingredients[input.IngredientIndex].ID,
	}

	return x
}

// ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionUpdateRequestInput creates a RecipeStepCompletionConditionUpdateRequestInput from a RecipeStepCompletionCondition.
func ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionUpdateRequestInput(input *types.RecipeStepCompletionCondition) *types.RecipeStepCompletionConditionUpdateRequestInput {
	x := &types.RecipeStepCompletionConditionUpdateRequestInput{
		IngredientStateID:   &input.IngredientState.ID,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Notes:               &input.Notes,
		Optional:            &input.Optional,
	}

	return x
}

// ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput builds a RecipeStepCompletionConditionCreationRequestInput from a RecipeStepCompletionCondition.
func ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput(recipeStepCompletionCondition *types.RecipeStepCompletionCondition, recipeStep *types.RecipeStep) *types.RecipeStepCompletionConditionCreationRequestInput {
	var ingredients []*types.RecipeStepCompletionConditionIngredientCreationRequestInput
	for _, i := range recipeStepCompletionCondition.Ingredients {
		x := ConvertRecipeStepCompletionConditionIngredientToRecipeStepCompletionConditionIngredientCreationRequestInput(i, recipeStep)
		ingredients = append(ingredients, x)
	}

	return &types.RecipeStepCompletionConditionCreationRequestInput{
		IngredientStateID:   recipeStepCompletionCondition.IngredientState.ID,
		BelongsToRecipeStep: recipeStepCompletionCondition.BelongsToRecipeStep,
		Notes:               recipeStepCompletionCondition.Notes,
		Ingredients:         ingredients,
		Optional:            recipeStepCompletionCondition.Optional,
	}
}

// ConvertRecipeStepCompletionConditionIngredientToRecipeStepCompletionConditionIngredientCreationRequestInput builds a RecipeStepCompletionConditionIngredientCreationRequestInput from a RecipeStepCompletionCondition.
func ConvertRecipeStepCompletionConditionIngredientToRecipeStepCompletionConditionIngredientCreationRequestInput(recipeStepCompletionConditionIngredient *types.RecipeStepCompletionConditionIngredient, recipeStep *types.RecipeStep) *types.RecipeStepCompletionConditionIngredientCreationRequestInput {
	var ingredientIndex uint64
	for i, ingredient := range recipeStep.Ingredients {
		if ingredient.ID == recipeStepCompletionConditionIngredient.RecipeStepIngredient {
			ingredientIndex = uint64(i)
		}
	}

	return &types.RecipeStepCompletionConditionIngredientCreationRequestInput{
		IngredientIndex: ingredientIndex,
	}
}

// ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput builds a RecipeStepCompletionConditionDatabaseCreationInput from a RecipeStepCompletionCondition.
func ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(recipeStepCompletionCondition *types.RecipeStepCompletionCondition) *types.RecipeStepCompletionConditionDatabaseCreationInput {
	return &types.RecipeStepCompletionConditionDatabaseCreationInput{
		ID:                  recipeStepCompletionCondition.ID,
		Optional:            recipeStepCompletionCondition.Optional,
		Notes:               recipeStepCompletionCondition.Notes,
		IngredientStateID:   recipeStepCompletionCondition.IngredientState.ID,
		BelongsToRecipeStep: recipeStepCompletionCondition.BelongsToRecipeStep,
	}
}
