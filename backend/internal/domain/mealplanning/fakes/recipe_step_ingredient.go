package fakes

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"

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
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepIngredient())
	}

	return &filtering.QueryFilteredResult[types.RecipeStepIngredient]{
		Pagination: filtering.Pagination{
			Cursor:          BuildFakeID(),
			MaxResponseSize: 50,
			FilteredCount:   exampleQuantity / 2,
			TotalCount:      exampleQuantity,
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
// Note: This now includes bridge table IDs since they are required.
func BuildFakeRecipeStepIngredientCreationRequestInput() *types.RecipeStepIngredientCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	input := converters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(recipeStepIngredient)
	// Bridge table IDs are now required
	input.ValidIngredientPreparationID = pointer.To(BuildFakeID())
	input.ValidIngredientMeasurementUnitID = pointer.To(BuildFakeID())
	return input
}

// BuildFakeRecipeStepIngredientCreationRequestInputForRecipeStepProduct builds a faked RecipeStepIngredientCreationRequestInput
// for a recipe step product (no bridge table IDs required).
func BuildFakeRecipeStepIngredientCreationRequestInputForRecipeStepProduct() *types.RecipeStepIngredientCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	input := converters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(recipeStepIngredient)
	input.ProductOfRecipeStepIndex = pointer.To(uint64(0))
	input.ProductOfRecipeStepProductIndex = pointer.To(uint64(0))
	return input
}
