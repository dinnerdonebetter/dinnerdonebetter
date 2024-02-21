package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepIngredient() *types.RecipeStepIngredient {
	return &types.RecipeStepIngredient{
		ID:                     BuildFakeID(),
		Name:                   buildUniqueString(),
		Ingredient:             BuildFakeValidIngredient(),
		MeasurementUnit:        *BuildFakeValidMeasurementUnit(),
		MinimumQuantity:        float32(BuildFakeNumber()),
		MaximumQuantity:        pointers.Pointer(float32(BuildFakeNumber())),
		QuantityNotes:          buildUniqueString(),
		Optional:               fake.Bool(),
		IngredientNotes:        buildUniqueString(),
		CreatedAt:              BuildFakeTime(),
		BelongsToRecipeStep:    BuildFakeID(),
		VesselIndex:            pointers.Pointer(fake.Uint16()),
		ToTaste:                fake.Bool(),
		ProductPercentageToUse: pointers.Pointer(float32(BuildFakeNumber())),
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
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepIngredientUpdateRequestInput builds a faked RecipeStepIngredientUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateRequestInput() *types.RecipeStepIngredientUpdateRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return converters.ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(recipeStepIngredient)
}

// BuildFakeRecipeStepIngredientCreationRequestInput builds a faked RecipeStepIngredientCreationRequestInput.
func BuildFakeRecipeStepIngredientCreationRequestInput() *types.RecipeStepIngredientCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return converters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(recipeStepIngredient)
}
