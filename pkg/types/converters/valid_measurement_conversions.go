package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertValidMeasurementConversionToValidMeasurementConversionUpdateRequestInput creates a ValidMeasurementUnitConversionUpdateRequestInput from a ValidMeasurementUnitConversion.
func ConvertValidMeasurementConversionToValidMeasurementConversionUpdateRequestInput(input *types.ValidMeasurementUnitConversion) *types.ValidMeasurementUnitConversionUpdateRequestInput {
	x := &types.ValidMeasurementUnitConversionUpdateRequestInput{
		From:     &input.From.ID,
		To:       &input.To.ID,
		Modifier: &input.Modifier,
		Notes:    &input.Notes,
	}

	if input.OnlyForIngredient != nil {
		x.OnlyForIngredient = &input.OnlyForIngredient.ID
	}

	return x
}

// ConvertValidMeasurementConversionCreationRequestInputToValidMeasurementConversionDatabaseCreationInput creates a ValidMeasurementConversionDatabaseCreationInput from a ValidMeasurementUnitConversionCreationRequestInput.
func ConvertValidMeasurementConversionCreationRequestInputToValidMeasurementConversionDatabaseCreationInput(input *types.ValidMeasurementUnitConversionCreationRequestInput) *types.ValidMeasurementConversionDatabaseCreationInput {
	x := &types.ValidMeasurementConversionDatabaseCreationInput{
		ID:                identifiers.New(),
		From:              input.From,
		To:                input.To,
		OnlyForIngredient: input.OnlyForIngredient,
		Modifier:          input.Modifier,
		Notes:             input.Notes,
	}

	return x
}

// ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput builds a ValidMeasurementUnitConversionCreationRequestInput from a ValidMeasurementUnitConversion.
func ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(validMeasurementConversion *types.ValidMeasurementUnitConversion) *types.ValidMeasurementUnitConversionCreationRequestInput {
	x := &types.ValidMeasurementUnitConversionCreationRequestInput{
		From:     validMeasurementConversion.From.ID,
		To:       validMeasurementConversion.To.ID,
		Modifier: validMeasurementConversion.Modifier,
		Notes:    validMeasurementConversion.Notes,
	}

	if validMeasurementConversion.OnlyForIngredient != nil {
		x.OnlyForIngredient = &validMeasurementConversion.OnlyForIngredient.ID
	}

	return x
}

// ConvertValidMeasurementConversionToValidMeasurementConversionDatabaseCreationInput builds a ValidMeasurementConversionDatabaseCreationInput from a ValidMeasurementUnitConversion.
func ConvertValidMeasurementConversionToValidMeasurementConversionDatabaseCreationInput(validMeasurementConversion *types.ValidMeasurementUnitConversion) *types.ValidMeasurementConversionDatabaseCreationInput {
	x := &types.ValidMeasurementConversionDatabaseCreationInput{
		ID:       validMeasurementConversion.ID,
		From:     validMeasurementConversion.From.ID,
		To:       validMeasurementConversion.To.ID,
		Modifier: validMeasurementConversion.Modifier,
		Notes:    validMeasurementConversion.Notes,
	}

	if validMeasurementConversion.OnlyForIngredient != nil {
		x.OnlyForIngredient = &validMeasurementConversion.OnlyForIngredient.ID
	}

	return x
}
