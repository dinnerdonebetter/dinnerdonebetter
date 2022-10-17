package converters

import "github.com/prixfixeco/api_server/pkg/types"

// ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput creates a RecipeStepProductUpdateRequestInput from a RecipeStepProduct.
func ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput(input *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
	if input == nil {
		return nil
	}

	x := &types.RecipeStepProductUpdateRequestInput{
		Name:                               &input.Name,
		Type:                               &input.Type,
		MeasurementUnitID:                  &input.MeasurementUnit.ID,
		QuantityNotes:                      &input.QuantityNotes,
		BelongsToRecipeStep:                &input.BelongsToRecipeStep,
		MinimumQuantity:                    &input.MinimumQuantity,
		MaximumQuantity:                    &input.MaximumQuantity,
		Compostable:                        &input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                &input.StorageInstructions,
	}

	return x
}

// ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput creates a RecipeStepProductDatabaseCreationInput from a RecipeStepProductCreationRequestInput.
func ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(input *types.RecipeStepProductCreationRequestInput) *types.RecipeStepProductDatabaseCreationInput {
	if input == nil {
		return nil
	}

	x := &types.RecipeStepProductDatabaseCreationInput{
		Name:                               input.Name,
		Type:                               input.Type,
		MeasurementUnitID:                  input.MeasurementUnitID,
		QuantityNotes:                      input.QuantityNotes,
		MinimumQuantity:                    input.MinimumQuantity,
		MaximumQuantity:                    input.MaximumQuantity,
		Compostable:                        input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                input.StorageInstructions,
	}

	return x
}

// ConvertRecipeStepProductToRecipeStepProductCreationRequestInput builds a RecipeStepProductCreationRequestInput from a RecipeStepProduct.
func ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductCreationRequestInput {
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

// ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput builds a RecipeStepProductDatabaseCreationInput from a RecipeStepProduct.
func ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductDatabaseCreationInput {
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