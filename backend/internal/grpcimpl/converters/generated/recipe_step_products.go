package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepProductCreationRequestInputToRecipeStepProduct(input *messages.RecipeStepProductCreationRequestInput) *messages.RecipeStepProduct {

output := &messages.RecipeStepProduct{
    IsLiquid: input.IsLiquid,
    StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
    Name: input.Name,
    StorageInstructions: input.StorageInstructions,
    QuantityNotes: input.QuantityNotes,
    Index: input.Index,
    Compostable: input.Compostable,
    StorageDurationInSeconds: input.StorageDurationInSeconds,
    Quantity: input.Quantity,
    Type: input.Type,
    ContainedInVesselIndex: input.ContainedInVesselIndex,
    IsWaste: input.IsWaste,
}

return output
}
func ConvertRecipeStepProductUpdateRequestInputToRecipeStepProduct(input *messages.RecipeStepProductUpdateRequestInput) *messages.RecipeStepProduct {

output := &messages.RecipeStepProduct{
    Name: input.Name,
    QuantityNotes: input.QuantityNotes,
    IsWaste: input.IsWaste,
    IsLiquid: input.IsLiquid,
    StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
    StorageDurationInSeconds: input.StorageDurationInSeconds,
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    ContainedInVesselIndex: input.ContainedInVesselIndex,
    Index: input.Index,
    Compostable: input.Compostable,
    Quantity: input.Quantity,
    Type: input.Type,
    StorageInstructions: input.StorageInstructions,
}

return output
}
