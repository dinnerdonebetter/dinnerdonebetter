package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	minQty := BuildFakeNumber()
	storageTemp := BuildFakeNumber()

	return &types.RecipeStepProduct{
		ID:                                 BuildFakeID(),
		Name:                               buildUniqueString(),
		Type:                               types.RecipeStepProductIngredientType,
		MinimumQuantity:                    float32(minQty),
		MaximumQuantity:                    float32(minQty + 1),
		QuantityNotes:                      buildUniqueString(),
		MeasurementUnit:                    *BuildFakeValidMeasurementUnit(),
		CreatedAt:                          fake.Date(),
		BelongsToRecipeStep:                fake.UUID(),
		Compostable:                        fake.Bool(),
		MaximumStorageDurationInSeconds:    pointers.Uint32Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius: pointers.Float32Pointer(float32(storageTemp)),
		MaximumStorageTemperatureInCelsius: pointers.Float32Pointer(float32(storageTemp + 1)),
		StorageInstructions:                buildUniqueString(),
	}
}

// BuildFakeRecipeStepProductList builds a faked RecipeStepProductList.
func BuildFakeRecipeStepProductList() *types.RecipeStepProductList {
	var examples []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepProduct())
	}

	return &types.RecipeStepProductList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeStepProducts: examples,
	}
}

// BuildFakeRecipeStepProductUpdateRequestInput builds a faked RecipeStepProductUpdateRequestInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateRequestInput() *types.RecipeStepProductUpdateRequestInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return &types.RecipeStepProductUpdateRequestInput{
		Name:                               &recipeStepProduct.Name,
		Type:                               &recipeStepProduct.Type,
		MinimumQuantity:                    &recipeStepProduct.MinimumQuantity,
		MaximumQuantity:                    &recipeStepProduct.MaximumQuantity,
		QuantityNotes:                      &recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  &recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                &recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        &recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: recipeStepProduct.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                &recipeStepProduct.StorageInstructions,
	}
}

// BuildFakeRecipeStepProductCreationRequestInput builds a faked RecipeStepProductCreationRequestInput.
func BuildFakeRecipeStepProductCreationRequestInput() *types.RecipeStepProductCreationRequestInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(recipeStepProduct)
}
