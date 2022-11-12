package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeValidIngredientMeasurementUnit builds a faked valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnit() *types.ValidIngredientMeasurementUnit {
	minQty := BuildFakeNumber()

	return &types.ValidIngredientMeasurementUnit{
		ID:                       BuildFakeID(),
		Notes:                    buildUniqueString(),
		MeasurementUnit:          *BuildFakeValidMeasurementUnit(),
		Ingredient:               *BuildFakeValidIngredient(),
		MinimumAllowableQuantity: float32(minQty),
		MaximumAllowableQuantity: float32(minQty + 1),
		CreatedAt:                BuildFakeTime(),
	}
}

// BuildFakeValidIngredientMeasurementUnitList builds a faked ValidIngredientMeasurementUnitList.
func BuildFakeValidIngredientMeasurementUnitList() *types.ValidIngredientMeasurementUnitList {
	var examples []*types.ValidIngredientMeasurementUnit
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientMeasurementUnit())
	}

	return &types.ValidIngredientMeasurementUnitList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidIngredientMeasurementUnits: examples,
	}
}

// BuildFakeValidIngredientMeasurementUnitUpdateRequestInput builds a faked ValidIngredientMeasurementUnitUpdateRequestInput from a valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnitUpdateRequestInput() *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                    &validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID:   &validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:        &validIngredientMeasurementUnit.Ingredient.ID,
		MinimumAllowableQuantity: &validIngredientMeasurementUnit.MinimumAllowableQuantity,
		MaximumAllowableQuantity: &validIngredientMeasurementUnit.MaximumAllowableQuantity,
	}
}

// BuildFakeValidIngredientMeasurementUnitCreationRequestInput builds a faked ValidIngredientMeasurementUnitCreationRequestInput.
func BuildFakeValidIngredientMeasurementUnitCreationRequestInput() *types.ValidIngredientMeasurementUnitCreationRequestInput {
	validIngredientMeasurementUnit := BuildFakeValidIngredientMeasurementUnit()
	return converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(validIngredientMeasurementUnit)
}
