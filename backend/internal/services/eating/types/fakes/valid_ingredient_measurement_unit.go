package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

// BuildFakeValidIngredientMeasurementUnit builds a faked valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnit() *types.ValidIngredientMeasurementUnit {
	return &types.ValidIngredientMeasurementUnit{
		ID:                BuildFakeID(),
		Notes:             buildUniqueString(),
		MeasurementUnit:   *BuildFakeValidMeasurementUnit(),
		Ingredient:        *BuildFakeValidIngredient(),
		AllowableQuantity: BuildFakeFloat32RangeWithOptionalMax(),
		CreatedAt:         BuildFakeTime(),
	}
}

// BuildFakeValidIngredientMeasurementUnitsList builds a faked ValidIngredientMeasurementUnitList.
func BuildFakeValidIngredientMeasurementUnitsList() *filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit] {
	var examples []*types.ValidIngredientMeasurementUnit
	for range exampleQuantity {
		examples = append(examples, BuildFakeValidIngredientMeasurementUnit())
	}

	return &filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientMeasurementUnitUpdateRequestInput builds a faked ValidIngredientMeasurementUnitUpdateRequestInput from a valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnitUpdateRequestInput() *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                  &validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID: &validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:      &validIngredientMeasurementUnit.Ingredient.ID,
		AllowableQuantity: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Min: &validIngredientMeasurementUnit.AllowableQuantity.Min,
			Max: validIngredientMeasurementUnit.AllowableQuantity.Max,
		},
	}
}

// BuildFakeValidIngredientMeasurementUnitCreationRequestInput builds a faked ValidIngredientMeasurementUnitCreationRequestInput.
func BuildFakeValidIngredientMeasurementUnitCreationRequestInput() *types.ValidIngredientMeasurementUnitCreationRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit)
}
