package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeValidIngredientMeasurementUnit builds a faked valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnit() *types.ValidIngredientMeasurementUnit {
	return &types.ValidIngredientMeasurementUnit{
		ID:                       ksuid.New().String(),
		Notes:                    buildUniqueString(),
		MeasurementUnit:          *BuildFakeValidMeasurementUnit(),
		Ingredient:               *BuildFakeValidIngredient(),
		MinimumAllowableQuantity: fake.Float32(),
		MaximumAllowableQuantity: fake.Float32(),
		CreatedAt:                uint64(uint32(fake.Date().Unix())),
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

// BuildFakeValidIngredientMeasurementUnitUpdateRequestInputFromValidIngredientMeasurementUnit builds a faked ValidIngredientMeasurementUnitUpdateRequestInput from a valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnitUpdateRequestInputFromValidIngredientMeasurementUnit(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitUpdateRequestInput {
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
	return BuildFakeValidIngredientMeasurementUnitCreationRequestInputFromValidIngredientMeasurementUnit(validIngredientMeasurementUnit)
}

// BuildFakeValidIngredientMeasurementUnitCreationRequestInputFromValidIngredientMeasurementUnit builds a faked ValidIngredientMeasurementUnitCreationRequestInput from a valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnitCreationRequestInputFromValidIngredientMeasurementUnit(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitCreationRequestInput {
	return &types.ValidIngredientMeasurementUnitCreationRequestInput{
		ID:                       validIngredientMeasurementUnit.ID,
		Notes:                    validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID:   validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:        validIngredientMeasurementUnit.Ingredient.ID,
		MinimumAllowableQuantity: validIngredientMeasurementUnit.MinimumAllowableQuantity,
		MaximumAllowableQuantity: validIngredientMeasurementUnit.MaximumAllowableQuantity,
	}
}

// BuildFakeValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnit builds a faked ValidIngredientMeasurementUnitDatabaseCreationInput from a valid ingredient measurement unit.
func BuildFakeValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnit(validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitDatabaseCreationInput {
	return &types.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                       validIngredientMeasurementUnit.ID,
		Notes:                    validIngredientMeasurementUnit.Notes,
		ValidMeasurementUnitID:   validIngredientMeasurementUnit.MeasurementUnit.ID,
		ValidIngredientID:        validIngredientMeasurementUnit.Ingredient.ID,
		MinimumAllowableQuantity: validIngredientMeasurementUnit.MinimumAllowableQuantity,
		MaximumAllowableQuantity: validIngredientMeasurementUnit.MaximumAllowableQuantity,
	}
}
