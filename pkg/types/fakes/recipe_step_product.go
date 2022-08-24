package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	return &types.RecipeStepProduct{
		ID:                                 ksuid.New().String(),
		Name:                               buildUniqueString(),
		Type:                               types.RecipeStepProductIngredientType,
		MinimumQuantityValue:               fake.Float32(),
		MaximumQuantityValue:               fake.Float32(),
		QuantityNotes:                      buildUniqueString(),
		MeasurementUnit:                    *BuildFakeValidMeasurementUnit(),
		CreatedOn:                          uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep:                fake.UUID(),
		Compostable:                        fake.Bool(),
		MaximumStorageDurationInSeconds:    fake.Uint32(),
		MinimumStorageTemperatureInCelsius: fake.Float32(),
		MaximumStorageTemperatureInCelsius: fake.Float32(),
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
		MinimumQuantityValue:               &recipeStepProduct.MinimumQuantityValue,
		MaximumQuantityValue:               &recipeStepProduct.MaximumQuantityValue,
		QuantityNotes:                      &recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  &recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                &recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        &recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    &recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: &recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: &recipeStepProduct.MaximumStorageTemperatureInCelsius,
	}
}

// BuildFakeRecipeStepProductUpdateRequestInputFromRecipeStepProduct builds a faked RecipeStepProductUpdateRequestInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateRequestInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
	return &types.RecipeStepProductUpdateRequestInput{
		Name:                               &recipeStepProduct.Name,
		Type:                               &recipeStepProduct.Type,
		MinimumQuantityValue:               &recipeStepProduct.MinimumQuantityValue,
		MaximumQuantityValue:               &recipeStepProduct.MaximumQuantityValue,
		QuantityNotes:                      &recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  &recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                &recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        &recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    &recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: &recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: &recipeStepProduct.MaximumStorageTemperatureInCelsius,
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
		MinimumQuantityValue:               recipeStepProduct.MinimumQuantityValue,
		MaximumQuantityValue:               recipeStepProduct.MaximumQuantityValue,
		QuantityNotes:                      recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: recipeStepProduct.MaximumStorageTemperatureInCelsius,
	}
}

// BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct builds a faked RecipeStepProductDatabaseCreationInput from a recipe step product.
func BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductDatabaseCreationInput {
	return &types.RecipeStepProductDatabaseCreationInput{
		ID:                                 recipeStepProduct.ID,
		Name:                               recipeStepProduct.Name,
		Type:                               recipeStepProduct.Type,
		MinimumQuantityValue:               recipeStepProduct.MinimumQuantityValue,
		MaximumQuantityValue:               recipeStepProduct.MaximumQuantityValue,
		QuantityNotes:                      recipeStepProduct.QuantityNotes,
		MeasurementUnitID:                  recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:                recipeStepProduct.BelongsToRecipeStep,
		Compostable:                        recipeStepProduct.Compostable,
		MaximumStorageDurationInSeconds:    recipeStepProduct.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: recipeStepProduct.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: recipeStepProduct.MaximumStorageTemperatureInCelsius,
	}
}
