package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertRecipeStepConditionCreationRequestInputToRecipeStepConditionDatabaseCreationInput creates a RecipeStepConditionDatabaseCreationInput from a RecipeStepConditionCreationRequestInput.
func ConvertRecipeStepConditionCreationRequestInputToRecipeStepConditionDatabaseCreationInput(input *types.RecipeStepConditionCreationRequestInput) *types.RecipeStepConditionDatabaseCreationInput {
	id := identifiers.New()

	var ingredients []*types.RecipeStepConditionIngredientDatabaseCreationInput
	for _, i := range input.Ingredients {
		x := ConvertRecipeStepConditionIngredientCreationRequestInputToRecipeStepConditionIngredientDatabaseCreationInput(i)
		x.BelongsToRecipeStepCondition = id
		ingredients = append(ingredients, x)
	}

	x := &types.RecipeStepConditionDatabaseCreationInput{
		ID:                  id,
		IngredientStateID:   input.IngredientStateID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Ingredients:         ingredients,
		Optional:            input.Optional,
	}

	return x
}

// ConvertRecipeStepConditionIngredientCreationRequestInputToRecipeStepConditionIngredientDatabaseCreationInput creates a RecipeStepConditionIngredientDatabaseCreationInput from a RecipeStepConditionCreationRequestInput.
func ConvertRecipeStepConditionIngredientCreationRequestInputToRecipeStepConditionIngredientDatabaseCreationInput(input *types.RecipeStepConditionIngredientCreationRequestInput) *types.RecipeStepConditionIngredientDatabaseCreationInput {
	x := &types.RecipeStepConditionIngredientDatabaseCreationInput{
		ID:                           identifiers.New(),
		BelongsToRecipeStepCondition: input.BelongsToRecipeStepCondition,
		RecipeStepIngredient:         input.RecipeStepIngredient,
	}

	return x
}

// ConvertRecipeStepConditionToRecipeStepConditionUpdateRequestInput creates a RecipeStepConditionUpdateRequestInput from a RecipeStepCondition.
func ConvertRecipeStepConditionToRecipeStepConditionUpdateRequestInput(input *types.RecipeStepCondition) *types.RecipeStepConditionUpdateRequestInput {
	x := &types.RecipeStepConditionUpdateRequestInput{
		IngredientStateID:   &input.IngredientState.ID,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Notes:               &input.Notes,
		Optional:            &input.Optional,
	}

	return x
}

// ConvertRecipeStepConditionToRecipeStepConditionCreationRequestInput builds a RecipeStepConditionCreationRequestInput from a RecipeStepCondition.
func ConvertRecipeStepConditionToRecipeStepConditionCreationRequestInput(recipeStepCondition *types.RecipeStepCondition) *types.RecipeStepConditionCreationRequestInput {
	var ingredients []*types.RecipeStepConditionIngredientCreationRequestInput
	for _, i := range recipeStepCondition.Ingredients {
		x := ConvertRecipeStepConditionIngredientToRecipeStepConditionIngredientCreationRequestInput(i)
		ingredients = append(ingredients, x)
	}

	return &types.RecipeStepConditionCreationRequestInput{
		IngredientStateID:   recipeStepCondition.IngredientState.ID,
		BelongsToRecipeStep: recipeStepCondition.BelongsToRecipeStep,
		Notes:               recipeStepCondition.Notes,
		Ingredients:         ingredients,
		Optional:            recipeStepCondition.Optional,
	}
}

// ConvertRecipeStepConditionIngredientToRecipeStepConditionIngredientCreationRequestInput builds a RecipeStepConditionIngredientCreationRequestInput from a RecipeStepCondition.
func ConvertRecipeStepConditionIngredientToRecipeStepConditionIngredientCreationRequestInput(recipeStepConditionIngredient *types.RecipeStepConditionIngredient) *types.RecipeStepConditionIngredientCreationRequestInput {
	return &types.RecipeStepConditionIngredientCreationRequestInput{
		BelongsToRecipeStepCondition: recipeStepConditionIngredient.BelongsToRecipeStepCondition,
		RecipeStepIngredient:         recipeStepConditionIngredient.RecipeStepIngredient,
	}
}

// ConvertRecipeStepConditionToRecipeStepConditionDatabaseCreationInput builds a RecipeStepConditionDatabaseCreationInput from a RecipeStepCondition.
func ConvertRecipeStepConditionToRecipeStepConditionDatabaseCreationInput(recipeStepCondition *types.RecipeStepCondition) *types.RecipeStepConditionDatabaseCreationInput {
	return &types.RecipeStepConditionDatabaseCreationInput{
		ID:                  recipeStepCondition.ID,
		Optional:            recipeStepCondition.Optional,
		BelongsToRecipeStep: recipeStepCondition.BelongsToRecipeStep,
	}
}
