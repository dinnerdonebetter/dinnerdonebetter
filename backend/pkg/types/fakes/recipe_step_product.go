package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	p := &types.RecipeStepProduct{
		ID:                          BuildFakeID(),
		Name:                        buildUniqueString(),
		Type:                        types.RecipeStepProductIngredientType,
		QuantityNotes:               buildUniqueString(),
		MeasurementUnit:             BuildFakeValidMeasurementUnit(),
		CreatedAt:                   BuildFakeTime(),
		BelongsToRecipeStep:         fake.UUID(),
		Compostable:                 fake.Bool(),
		IsLiquid:                    fake.Bool(),
		IsWaste:                     fake.Bool(),
		Quantity:                    BuildFakeOptionalFloat32Range(),
		StorageDurationInSeconds:    BuildFakeOptionalUint32Range(),
		StorageTemperatureInCelsius: BuildFakeOptionalFloat32Range(),
		StorageInstructions:         buildUniqueString(),
		Index:                       fake.Uint16(),
		ContainedInVesselIndex:      pointer.To(fake.Uint16()),
	}

	// TODO: there's no database field for this
	p.StorageDurationInSeconds.Min = nil

	return p
}

// BuildFakeRecipeStepProductsList builds a faked RecipeStepProductList.
func BuildFakeRecipeStepProductsList() *filtering.QueryFilteredResult[types.RecipeStepProduct] {
	var examples []*types.RecipeStepProduct
	for range exampleQuantity {
		examples = append(examples, BuildFakeRecipeStepProduct())
	}

	return &filtering.QueryFilteredResult[types.RecipeStepProduct]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepProductUpdateRequestInput builds a faked RecipeStepProductUpdateRequestInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateRequestInput() *types.RecipeStepProductUpdateRequestInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return converters.ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput(recipeStepProduct)
}

// BuildFakeRecipeStepProductCreationRequestInput builds a faked RecipeStepProductCreationRequestInput.
func BuildFakeRecipeStepProductCreationRequestInput() *types.RecipeStepProductCreationRequestInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(recipeStepProduct)
}
