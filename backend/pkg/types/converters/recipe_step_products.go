package converters

import (
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput creates a RecipeStepProductUpdateRequestInput from a RecipeStepProduct.
func ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput(input *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
	if input == nil {
		return nil
	}

	x := &types.RecipeStepProductUpdateRequestInput{
		Name:                        &input.Name,
		Type:                        &input.Type,
		MeasurementUnitID:           &input.MeasurementUnit.ID,
		QuantityNotes:               &input.QuantityNotes,
		BelongsToRecipeStep:         &input.BelongsToRecipeStep,
		Compostable:                 &input.Compostable,
		Quantity:                    input.Quantity,
		StorageDurationInSeconds:    input.StorageDurationInSeconds,
		StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
		StorageInstructions:         &input.StorageInstructions,
		IsWaste:                     &input.IsWaste,
		IsLiquid:                    &input.IsLiquid,
		Index:                       &input.Index,
		ContainedInVesselIndex:      input.ContainedInVesselIndex,
	}

	return x
}

// ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput creates a RecipeStepProductDatabaseCreationInput from a RecipeStepProductCreationRequestInput.
func ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(input *types.RecipeStepProductCreationRequestInput) *types.RecipeStepProductDatabaseCreationInput {
	if input == nil {
		return nil
	}

	x := &types.RecipeStepProductDatabaseCreationInput{
		ID:                          identifiers.New(),
		Name:                        input.Name,
		Type:                        input.Type,
		MeasurementUnitID:           input.MeasurementUnitID,
		QuantityNotes:               input.QuantityNotes,
		Compostable:                 input.Compostable,
		Quantity:                    input.Quantity,
		StorageDurationInSeconds:    input.StorageDurationInSeconds,
		StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
		StorageInstructions:         input.StorageInstructions,
		IsWaste:                     input.IsWaste,
		IsLiquid:                    input.IsLiquid,
		Index:                       input.Index,
		ContainedInVesselIndex:      input.ContainedInVesselIndex,
	}

	return x
}

// ConvertRecipeStepProductToRecipeStepProductCreationRequestInput builds a RecipeStepProductCreationRequestInput from a RecipeStepProduct.
func ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductCreationRequestInput {
	return &types.RecipeStepProductCreationRequestInput{
		Name:                        recipeStepProduct.Name,
		Type:                        recipeStepProduct.Type,
		QuantityNotes:               recipeStepProduct.QuantityNotes,
		MeasurementUnitID:           &recipeStepProduct.MeasurementUnit.ID,
		Compostable:                 recipeStepProduct.Compostable,
		Quantity:                    recipeStepProduct.Quantity,
		StorageDurationInSeconds:    recipeStepProduct.StorageDurationInSeconds,
		StorageTemperatureInCelsius: recipeStepProduct.StorageTemperatureInCelsius,
		StorageInstructions:         recipeStepProduct.StorageInstructions,
		IsWaste:                     recipeStepProduct.IsWaste,
		IsLiquid:                    recipeStepProduct.IsLiquid,
		Index:                       recipeStepProduct.Index,
		ContainedInVesselIndex:      recipeStepProduct.ContainedInVesselIndex,
	}
}

// ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput builds a RecipeStepProductDatabaseCreationInput from a RecipeStepProduct.
func ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductDatabaseCreationInput {
	return &types.RecipeStepProductDatabaseCreationInput{
		ID:                          recipeStepProduct.ID,
		Name:                        recipeStepProduct.Name,
		Type:                        recipeStepProduct.Type,
		QuantityNotes:               recipeStepProduct.QuantityNotes,
		MeasurementUnitID:           &recipeStepProduct.MeasurementUnit.ID,
		BelongsToRecipeStep:         recipeStepProduct.BelongsToRecipeStep,
		Compostable:                 recipeStepProduct.Compostable,
		Quantity:                    recipeStepProduct.Quantity,
		StorageDurationInSeconds:    recipeStepProduct.StorageDurationInSeconds,
		StorageTemperatureInCelsius: recipeStepProduct.StorageTemperatureInCelsius,
		StorageInstructions:         recipeStepProduct.StorageInstructions,
		IsWaste:                     recipeStepProduct.IsWaste,
		IsLiquid:                    recipeStepProduct.IsLiquid,
		Index:                       recipeStepProduct.Index,
		ContainedInVesselIndex:      recipeStepProduct.ContainedInVesselIndex,
	}
}
