package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepCompletionCondition builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepCompletionCondition() *types.RecipeStepCompletionCondition {
	var ingredients []*types.RecipeStepCompletionConditionIngredient
	for i := 0; i < exampleQuantity; i++ {
		ingredients = append(ingredients, BuildFakeRecipeStepCompletionConditionIngredient())
	}

	return &types.RecipeStepCompletionCondition{
		Optional:            fake.Bool(),
		IngredientState:     *BuildFakeValidIngredientState(),
		ID:                  BuildFakeID(),
		BelongsToRecipeStep: BuildFakeID(),
		Notes:               buildUniqueString(),
		Ingredients:         ingredients,
	}
}

// BuildFakeRecipeStepCompletionConditionIngredient builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepCompletionConditionIngredient() *types.RecipeStepCompletionConditionIngredient {
	return &types.RecipeStepCompletionConditionIngredient{
		ID:                                     BuildFakeID(),
		BelongsToRecipeStepCompletionCondition: BuildFakeID(),
		RecipeStepIngredient:                   BuildFakeID(),
	}
}

// BuildFakeRecipeStepCompletionConditionList builds a faked RecipeStepCompletionConditionList.
func BuildFakeRecipeStepCompletionConditionList() *types.QueryFilteredResult[types.RecipeStepCompletionCondition] {
	var examples []*types.RecipeStepCompletionCondition
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepCompletionCondition())
	}

	return &types.QueryFilteredResult[types.RecipeStepCompletionCondition]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepCompletionConditionUpdateRequestInput builds a faked RecipeStepCompletionConditionUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepCompletionConditionUpdateRequestInput() *types.RecipeStepCompletionConditionUpdateRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepCompletionCondition()
	return &types.RecipeStepCompletionConditionUpdateRequestInput{
		Optional:            &recipeStepIngredient.Optional,
		BelongsToRecipeStep: &recipeStepIngredient.BelongsToRecipeStep,
		IngredientStateID:   &recipeStepIngredient.IngredientState.ID,
		Notes:               &recipeStepIngredient.Notes,
	}
}

// RecipeStepCompletionConditionForExistingRecipeCreationRequestInput builds a faked RecipeStepCompletionConditionForExistingRecipeCreationRequestInput.
func RecipeStepCompletionConditionForExistingRecipeCreationRequestInput() *types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepCompletionCondition()
	return converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionForExistingRecipeCreationRequestInput(recipeStepIngredient)
}
