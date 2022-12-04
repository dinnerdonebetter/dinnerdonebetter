package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeRecipeStepCondition builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepCondition() *types.RecipeStepCondition {
	var ingredients []*types.RecipeStepConditionIngredient
	for i := 0; i < exampleQuantity; i++ {
		ingredients = append(ingredients, BuildFakeRecipeStepConditionIngredient())
	}

	return &types.RecipeStepCondition{
		Optional:            fake.Bool(),
		IngredientState:     *BuildFakeValidIngredientState(),
		ID:                  BuildFakeID(),
		BelongsToRecipeStep: BuildFakeID(),
		Notes:               buildUniqueString(),
		Ingredients:         ingredients,
	}
}

// BuildFakeRecipeStepConditionIngredient builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepConditionIngredient() *types.RecipeStepConditionIngredient {
	return &types.RecipeStepConditionIngredient{
		ID:                           BuildFakeID(),
		BelongsToRecipeStepCondition: BuildFakeID(),
		RecipeStepIngredient:         BuildFakeID(),
	}
}

// BuildFakeRecipeStepConditionList builds a faked RecipeStepConditionList.
func BuildFakeRecipeStepConditionList() *types.QueryFilteredResult[types.RecipeStepCondition] {
	var examples []*types.RecipeStepCondition
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepCondition())
	}

	return &types.QueryFilteredResult[types.RecipeStepCondition]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepConditionUpdateRequestInput builds a faked RecipeStepConditionUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepConditionUpdateRequestInput() *types.RecipeStepConditionUpdateRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepCondition()
	return &types.RecipeStepConditionUpdateRequestInput{
		Optional:            &recipeStepIngredient.Optional,
		BelongsToRecipeStep: &recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepConditionCreationRequestInput builds a faked RecipeStepConditionCreationRequestInput.
func BuildFakeRecipeStepConditionCreationRequestInput() *types.RecipeStepConditionCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepCondition()
	return converters.ConvertRecipeStepConditionToRecipeStepConditionCreationRequestInput(recipeStepIngredient)
}
