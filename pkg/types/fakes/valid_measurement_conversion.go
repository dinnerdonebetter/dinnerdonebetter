package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeValidMeasurementConversion builds a faked valid preparation.
func BuildFakeValidMeasurementConversion() *types.ValidMeasurementUnitConversion {
	return &types.ValidMeasurementUnitConversion{
		ID:                BuildFakeID(),
		From:              *BuildFakeValidMeasurementUnit(),
		To:                *BuildFakeValidMeasurementUnit(),
		OnlyForIngredient: nil,
		Modifier:          float32(BuildFakeNumber()),
		Notes:             buildUniqueString(),
		CreatedAt:         BuildFakeTime(),
	}
}

// BuildFakeValidMeasurementConversionList builds a faked ValidMeasurementConversionList.
func BuildFakeValidMeasurementConversionList() *types.QueryFilteredResult[types.ValidMeasurementUnitConversion] {
	var examples []*types.ValidMeasurementUnitConversion
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidMeasurementConversion())
	}

	return &types.QueryFilteredResult[types.ValidMeasurementUnitConversion]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidMeasurementConversionUpdateRequestInput builds a faked ValidMeasurementUnitConversionUpdateRequestInput from a valid preparation.
func BuildFakeValidMeasurementConversionUpdateRequestInput() *types.ValidMeasurementUnitConversionUpdateRequestInput {
	validMeasurementConversion := BuildFakeValidMeasurementConversion()

	x := &types.ValidMeasurementUnitConversionUpdateRequestInput{
		From:     &validMeasurementConversion.From.ID,
		To:       &validMeasurementConversion.To.ID,
		Modifier: &validMeasurementConversion.Modifier,
		Notes:    &validMeasurementConversion.Notes,
	}

	if validMeasurementConversion.OnlyForIngredient != nil {
		x.OnlyForIngredient = &validMeasurementConversion.OnlyForIngredient.ID
	}

	return x
}

// BuildFakeValidMeasurementConversionCreationRequestInput builds a faked ValidMeasurementUnitConversionCreationRequestInput.
func BuildFakeValidMeasurementConversionCreationRequestInput() *types.ValidMeasurementUnitConversionCreationRequestInput {
	validMeasurementConversion := BuildFakeValidMeasurementConversion()
	return converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(validMeasurementConversion)
}
