package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	return &types.RecipeStepProduct{
		ID:                                 ksuid.New().String(),
		Name:                               buildUniqueString(),
		Type:                               types.RecipeStepProductIngredientType,
		MinimumQuantity:                    fake.Float32(),
		MaximumQuantity:                    fake.Float32(),
		QuantityNotes:                      buildUniqueString(),
		MeasurementUnit:                    *BuildFakeValidMeasurementUnit(),
		CreatedAt:                          fake.Date(),
		BelongsToRecipeStep:                fake.UUID(),
		Compostable:                        fake.Bool(),
		MaximumStorageDurationInSeconds:    pointers.Uint32Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius: pointers.Float32Pointer(fake.Float32()),
		MaximumStorageTemperatureInCelsius: pointers.Float32Pointer(fake.Float32()),
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

// BuildFakeRecipeStepProductUpdateRequestInputFromRecipeStepProduct builds a faked RecipeStepProductUpdateRequestInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateRequestInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
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
	return BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(recipeStepProduct)
}

// BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct builds a faked RecipeStepProductCreationRequestInput from a recipe step product.
func BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductCreationRequestInput {
	return &types.RecipeStepProductCreationRequestInput{
		ID:                                 recipeStepProduct.ID,
		Name:                               recipeStepProduct.Name,
		Type:                               recipeStepProduct.Type,
		MinimumQuantity:                    recipeStepProduct.MinimumQuantity,
		MaximumQuantity:                    recipeStepProduct.MaximumQuantity,
		QuantityNotes:                      recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: recipeStepProduct.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                recipeStepProduct.StorageInstructions,
	}
}

// BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct builds a faked RecipeStepProductDatabaseCreationInput from a recipe step product.
func BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductDatabaseCreationInput {
	return &types.RecipeStepProductDatabaseCreationInput{
		ID:                                 recipeStepProduct.ID,
		Name:                               recipeStepProduct.Name,
		Type:                               recipeStepProduct.Type,
		MinimumQuantity:                    recipeStepProduct.MinimumQuantity,
		MaximumQuantity:                    recipeStepProduct.MaximumQuantity,
		QuantityNotes:                      recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: recipeStepProduct.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                recipeStepProduct.StorageInstructions,
	}
}
