package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertValidMeasurementConversionToValidMeasurementConversionUpdateRequestInput creates a ValidMeasurementConversionUpdateRequestInput from a ValidMeasurementConversion.
func ConvertValidMeasurementConversionToValidMeasurementConversionUpdateRequestInput(input *types.ValidMeasurementConversion) *types.ValidMeasurementConversionUpdateRequestInput {
	x := &types.ValidMeasurementConversionUpdateRequestInput{
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

// ConvertValidMeasurementConversionCreationRequestInputToValidMeasurementConversionDatabaseCreationInput creates a ValidMeasurementConversionDatabaseCreationInput from a ValidMeasurementConversionCreationRequestInput.
func ConvertValidMeasurementConversionCreationRequestInputToValidMeasurementConversionDatabaseCreationInput(input *types.ValidMeasurementConversionCreationRequestInput) *types.ValidMeasurementConversionDatabaseCreationInput {
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

// ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput builds a ValidMeasurementConversionCreationRequestInput from a ValidMeasurementConversion.
func ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(validMeasurementConversion *types.ValidMeasurementConversion) *types.ValidMeasurementConversionCreationRequestInput {
	x := &types.ValidMeasurementConversionCreationRequestInput{
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

// ConvertValidMeasurementConversionToValidMeasurementConversionDatabaseCreationInput builds a ValidMeasurementConversionDatabaseCreationInput from a ValidMeasurementConversion.
func ConvertValidMeasurementConversionToValidMeasurementConversionDatabaseCreationInput(validMeasurementConversion *types.ValidMeasurementConversion) *types.ValidMeasurementConversionDatabaseCreationInput {
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
