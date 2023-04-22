package fakes

import (
	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	minQty := float32(BuildFakeNumber())
	storageTemp := BuildFakeNumber()

	return &types.RecipeStepProduct{
		ID:                                 BuildFakeID(),
		Name:                               buildUniqueString(),
		Type:                               types.RecipeStepProductIngredientType,
		MinimumQuantity:                    pointers.Pointer(float32(minQty)),
		MaximumQuantity:                    pointers.Pointer(float32(minQty + 1)),
		QuantityNotes:                      buildUniqueString(),
		MeasurementUnit:                    BuildFakeValidMeasurementUnit(),
		CreatedAt:                          BuildFakeTime(),
		BelongsToRecipeStep:                fake.UUID(),
		Compostable:                        fake.Bool(),
		IsLiquid:                           fake.Bool(),
		IsWaste:                            fake.Bool(),
		MaximumStorageDurationInSeconds:    pointers.Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius: pointers.Pointer(float32(storageTemp)),
		MaximumStorageTemperatureInCelsius: pointers.Pointer(float32(storageTemp + 1)),
		StorageInstructions:                buildUniqueString(),
		Index:                              fake.Uint16(),
		ContainedInVesselIndex:             pointers.Pointer(fake.Uint16()),
	}
}

// BuildFakeRecipeStepProductList builds a faked RecipeStepProductList.
func BuildFakeRecipeStepProductList() *types.QueryFilteredResult[types.RecipeStepProduct] {
	var examples []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepProduct())
	}

	return &types.QueryFilteredResult[types.RecipeStepProduct]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
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
