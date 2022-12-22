package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepIngredient() *types.RecipeStepIngredient {
	return &types.RecipeStepIngredient{
		ID:                  BuildFakeID(),
		Name:                buildUniqueString(),
		Ingredient:          BuildFakeValidIngredient(),
		MeasurementUnit:     *BuildFakeValidMeasurementUnit(),
		MinimumQuantity:     float32(fake.Uint32()),
		MaximumQuantity:     float32(fake.Uint32()),
		QuantityNotes:       buildUniqueString(),
		Optional:            fake.Bool(),
		IngredientNotes:     buildUniqueString(),
		CreatedAt:           BuildFakeTime(),
		BelongsToRecipeStep: BuildFakeID(),
	}
}

// BuildFakeRecipeStepIngredientList builds a faked RecipeStepIngredientList.
func BuildFakeRecipeStepIngredientList() *types.QueryFilteredResult[types.RecipeStepIngredient] {
	var examples []*types.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepIngredient())
	}

	return &types.QueryFilteredResult[types.RecipeStepIngredient]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepIngredientUpdateRequestInput builds a faked RecipeStepIngredientUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateRequestInput() *types.RecipeStepIngredientUpdateRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return &types.RecipeStepIngredientUpdateRequestInput{
		Name:                &recipeStepIngredient.Name,
		Optional:            &recipeStepIngredient.Optional,
		IngredientID:        &recipeStepIngredient.Ingredient.ID,
		MeasurementUnitID:   &recipeStepIngredient.MeasurementUnit.ID,
		MinimumQuantity:     &recipeStepIngredient.MinimumQuantity,
		MaximumQuantity:     &recipeStepIngredient.MaximumQuantity,
		QuantityNotes:       &recipeStepIngredient.QuantityNotes,
		IngredientNotes:     &recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: &recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientCreationRequestInput builds a faked RecipeStepIngredientCreationRequestInput.
func BuildFakeRecipeStepIngredientCreationRequestInput() *types.RecipeStepIngredientCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return converters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(recipeStepIngredient)
}
