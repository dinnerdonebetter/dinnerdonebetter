package converters

import (
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput creates a ValidIngredientMeasurementUnitUpdateRequestInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(input *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	x := &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                    &input.Notes,
		ValidMeasurementUnitID:   &input.MeasurementUnit.ID,
		ValidIngredientID:        &input.Ingredient.ID,
		MinimumAllowableQuantity: &input.MinimumAllowableQuantity,
		MaximumAllowableQuantity: &input.MaximumAllowableQuantity,
	}

	return x
}

// ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput creates a ValidIngredientMeasurementUnitDatabaseCreationInput from a ValidIngredientMeasurementUnitCreationRequestInput.
func ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput(input *types.ValidIngredientMeasurementUnitCreationRequestInput) *types.ValidIngredientMeasurementUnitDatabaseCreationInput {
	x := &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		Notes:                    input.Notes,
		ValidMeasurementUnitID:   input.ValidMeasurementUnitID,
		ValidIngredientID:        input.ValidIngredientID,
		MinimumAllowableQuantity: input.MinimumAllowableQuantity,
		MaximumAllowableQuantity: input.MaximumAllowableQuantity,
	}

	return x
}

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput builds a ValidIngredientMeasurementUnitCreationRequestInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitCreationRequestInput {
	return &types.ValidIngredientMeasurementUnitCreationRequestInput{
		ID:                       validIngredientMeasurementUnit.ID,
		Notes:                    validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID:   validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:        validIngredientMeasurementUnit.Ingredient.ID,
		MinimumAllowableQuantity: validIngredientMeasurementUnit.MinimumAllowableQuantity,
		MaximumAllowableQuantity: validIngredientMeasurementUnit.MaximumAllowableQuantity,
	}
}

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput builds a ValidIngredientMeasurementUnitDatabaseCreationInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitDatabaseCreationInput {
	return &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                       validIngredientMeasurementUnit.ID,
		Notes:                    validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID:   validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:        validIngredientMeasurementUnit.Ingredient.ID,
		MinimumAllowableQuantity: validIngredientMeasurementUnit.MinimumAllowableQuantity,
		MaximumAllowableQuantity: validIngredientMeasurementUnit.MaximumAllowableQuantity,
	}
}