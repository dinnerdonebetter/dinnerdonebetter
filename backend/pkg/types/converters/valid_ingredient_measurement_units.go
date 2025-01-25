package converters

import (
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput creates a ValidIngredientMeasurementUnitUpdateRequestInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(input *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	x := &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                  &input.Notes,
		ValidMeasurementUnitID: &input.MeasurementUnit.ID,
		ValidIngredientID:      &input.Ingredient.ID,
		AllowableQuantity: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Max: input.AllowableQuantity.Max,
			Min: &input.AllowableQuantity.Min,
		},
	}

	return x
}

// ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput creates a ValidIngredientMeasurementUnitDatabaseCreationInput from a ValidIngredientMeasurementUnitCreationRequestInput.
func ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput(input *types.ValidIngredientMeasurementUnitCreationRequestInput) *types.ValidIngredientMeasurementUnitDatabaseCreationInput {
	x := &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     identifiers.New(),
		Notes:                  input.Notes,
		ValidMeasurementUnitID: input.ValidMeasurementUnitID,
		ValidIngredientID:      input.ValidIngredientID,
		AllowableQuantity:      input.AllowableQuantity,
	}

	return x
}

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput builds a ValidIngredientMeasurementUnitCreationRequestInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitCreationRequestInput {
	return &types.ValidIngredientMeasurementUnitCreationRequestInput{
		Notes:                  validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID: validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:      validIngredientMeasurementUnit.Ingredient.ID,
		AllowableQuantity:      validIngredientMeasurementUnit.AllowableQuantity,
	}
}

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput builds a ValidIngredientMeasurementUnitDatabaseCreationInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitDatabaseCreationInput {
	return &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     validIngredientMeasurementUnit.ID,
		Notes:                  validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID: validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:      validIngredientMeasurementUnit.Ingredient.ID,
		AllowableQuantity:      validIngredientMeasurementUnit.AllowableQuantity,
	}
}
