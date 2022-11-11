package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeValidMeasurementConversion builds a faked valid preparation.
func BuildFakeValidMeasurementConversion() *types.ValidMeasurementConversion {
	return &types.ValidMeasurementConversion{
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
func BuildFakeValidMeasurementConversionList() *types.ValidMeasurementConversionList {
	var examples []*types.ValidMeasurementConversion
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidMeasurementConversion())
	}

	return &types.ValidMeasurementConversionList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidMeasurementConversions: examples,
	}
}

// BuildFakeValidMeasurementConversionUpdateRequestInput builds a faked ValidMeasurementConversionUpdateRequestInput from a valid preparation.
func BuildFakeValidMeasurementConversionUpdateRequestInput() *types.ValidMeasurementConversionUpdateRequestInput {
	validMeasurementConversion := BuildFakeValidMeasurementConversion()

	x := &types.ValidMeasurementConversionUpdateRequestInput{
		From:     &validMeasurementConversion.From.ID,
		To:       &validMeasurementConversion.To.ID,
		Modifier: &validMeasurementConversion.Modifier,
		Notes:    &validMeasurementConversion.Notes,
	}

	if validMeasurementConversion.OnlyForIngredient != nil {
		x.ForIngredient = &validMeasurementConversion.OnlyForIngredient.ID
	}

	return x
}

// BuildFakeValidMeasurementConversionCreationRequestInput builds a faked ValidMeasurementConversionCreationRequestInput.
func BuildFakeValidMeasurementConversionCreationRequestInput() *types.ValidMeasurementConversionCreationRequestInput {
	validMeasurementConversion := BuildFakeValidMeasurementConversion()
	return converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(validMeasurementConversion)
}
