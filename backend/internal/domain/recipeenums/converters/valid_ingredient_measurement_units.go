package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput creates a ValidIngredientMeasurementUnitUpdateRequestInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(input *recipeenums.ValidIngredientMeasurementUnit) *recipeenums.ValidIngredientMeasurementUnitUpdateRequestInput {
	x := &recipeenums.ValidIngredientMeasurementUnitUpdateRequestInput{
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
func ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnitDatabaseCreationInput(input *recipeenums.ValidIngredientMeasurementUnitCreationRequestInput) *recipeenums.ValidIngredientMeasurementUnitDatabaseCreationInput {
	x := &recipeenums.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     identifiers.New(),
		Notes:                  input.Notes,
		ValidMeasurementUnitID: input.ValidMeasurementUnitID,
		ValidIngredientID:      input.ValidIngredientID,
		AllowableQuantity:      input.AllowableQuantity,
	}

	return x
}

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput builds a ValidIngredientMeasurementUnitCreationRequestInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit *recipeenums.ValidIngredientMeasurementUnit) *recipeenums.ValidIngredientMeasurementUnitCreationRequestInput {
	return &recipeenums.ValidIngredientMeasurementUnitCreationRequestInput{
		Notes:                  validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID: validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:      validIngredientMeasurementUnit.Ingredient.ID,
		AllowableQuantity:      validIngredientMeasurementUnit.AllowableQuantity,
	}
}

// ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput builds a ValidIngredientMeasurementUnitDatabaseCreationInput from a ValidIngredientMeasurementUnit.
func ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput(validIngredientMeasurementUnit *recipeenums.ValidIngredientMeasurementUnit) *recipeenums.ValidIngredientMeasurementUnitDatabaseCreationInput {
	return &recipeenums.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     validIngredientMeasurementUnit.ID,
		Notes:                  validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID: validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:      validIngredientMeasurementUnit.Ingredient.ID,
		AllowableQuantity:      validIngredientMeasurementUnit.AllowableQuantity,
	}
}
