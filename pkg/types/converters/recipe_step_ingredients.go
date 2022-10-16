package converters

import "github.com/prixfixeco/api_server/pkg/types"

// ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(input *types.RecipeStepIngredientCreationRequestInput) *types.RecipeStepIngredientDatabaseCreationInput {
	x := &types.RecipeStepIngredientDatabaseCreationInput{
		IngredientID:        input.IngredientID,
		Name:                input.Name,
		MeasurementUnitID:   input.MeasurementUnitID,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		IngredientNotes:     input.IngredientNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Optional:            input.Optional,
	}

	return x
}
