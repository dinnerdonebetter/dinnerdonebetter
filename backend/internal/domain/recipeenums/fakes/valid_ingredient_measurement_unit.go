package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// BuildFakeValidIngredientMeasurementUnit builds a faked valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnit() *recipeenums.ValidIngredientMeasurementUnit {
	return &recipeenums.ValidIngredientMeasurementUnit{
		ID:                BuildFakeID(),
		Notes:             buildUniqueString(),
		MeasurementUnit:   *BuildFakeValidMeasurementUnit(),
		Ingredient:        *BuildFakeValidIngredient(),
		AllowableQuantity: BuildFakeFloat32RangeWithOptionalMax(),
		CreatedAt:         BuildFakeTime(),
	}
}

// BuildFakeValidIngredientMeasurementUnitsList builds a faked ValidIngredientMeasurementUnitList.
func BuildFakeValidIngredientMeasurementUnitsList() *filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit] {
	var examples []*recipeenums.ValidIngredientMeasurementUnit
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientMeasurementUnit())
	}

	return &filtering.QueryFilteredResult[recipeenums.ValidIngredientMeasurementUnit]{
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
func BuildFakeValidIngredientMeasurementUnitUpdateRequestInput() *recipeenums.ValidIngredientMeasurementUnitUpdateRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return &recipeenums.ValidIngredientMeasurementUnitUpdateRequestInput{
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
func BuildFakeValidIngredientMeasurementUnitCreationRequestInput() *recipeenums.ValidIngredientMeasurementUnitCreationRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit)
}
