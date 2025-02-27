package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
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
		Quantity:               BuildFakeFloat32RangeWithOptionalMax(),
		QuantityNotes:          buildUniqueString(),
		Optional:               fake.Bool(),
		IngredientNotes:        buildUniqueString(),
		CreatedAt:              BuildFakeTime(),
		BelongsToRecipeStep:    BuildFakeID(),
		VesselIndex:            pointer.To(fake.Uint16()),
		ToTaste:                fake.Bool(),
		ProductPercentageToUse: pointer.To(float32(buildFakeNumber())),
	}
}

// BuildFakeRecipeStepIngredientsList builds a faked RecipeStepIngredientList.
func BuildFakeRecipeStepIngredientsList() *filtering.QueryFilteredResult[types.RecipeStepIngredient] {
	var examples []*types.RecipeStepIngredient
	for range exampleQuantity {
		examples = append(examples, BuildFakeRecipeStepIngredient())
	}

	return &filtering.QueryFilteredResult[types.RecipeStepIngredient]{
		Pagination: filtering.Pagination{
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
