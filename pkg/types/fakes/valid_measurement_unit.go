package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeValidMeasurementUnit builds a faked valid ingredient.
func BuildFakeValidMeasurementUnit() *types.ValidMeasurementUnit {
	return &types.ValidMeasurementUnit{
		ID:          ksuid.New().String(),
		Name:        buildUniqueString(),
		Description: buildUniqueString(),
		Volumetric:  fake.Bool(),
		IconPath:    buildUniqueString(),
		Universal:   fake.Bool(),
		Metric:      fake.Bool(),
		Imperial:    fake.Bool(),
		PluralName:  buildUniqueString(),
		CreatedOn:   uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidMeasurementUnitList builds a faked ValidMeasurementUnitList.
func BuildFakeValidMeasurementUnitList() *types.ValidMeasurementUnitList {
	var examples []*types.ValidMeasurementUnit
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidMeasurementUnit())
	}

	return &types.ValidMeasurementUnitList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidMeasurementUnits: examples,
	}
}

// BuildFakeValidMeasurementUnitUpdateRequestInput builds a faked ValidMeasurementUnitUpdateRequestInput from a valid ingredient.
func BuildFakeValidMeasurementUnitUpdateRequestInput() *types.ValidMeasurementUnitUpdateRequestInput {
	validMeasurementUnit := BuildFakeValidMeasurementUnit()
	return &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        &validMeasurementUnit.Name,
		Description: &validMeasurementUnit.Description,
		Volumetric:  &validMeasurementUnit.Volumetric,
		IconPath:    &validMeasurementUnit.IconPath,
		Universal:   &validMeasurementUnit.Universal,
		Metric:      &validMeasurementUnit.Metric,
		Imperial:    &validMeasurementUnit.Imperial,
		PluralName:  &validMeasurementUnit.PluralName,
	}
}

// BuildFakeValidMeasurementUnitUpdateRequestInputFromValidMeasurementUnit builds a faked ValidMeasurementUnitUpdateRequestInput from a valid ingredient.
func BuildFakeValidMeasurementUnitUpdateRequestInputFromValidMeasurementUnit(validMeasurementUnit *types.ValidMeasurementUnit) *types.ValidMeasurementUnitUpdateRequestInput {
	return &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        &validMeasurementUnit.Name,
		Description: &validMeasurementUnit.Description,
		Volumetric:  &validMeasurementUnit.Volumetric,
		IconPath:    &validMeasurementUnit.IconPath,
		Universal:   &validMeasurementUnit.Universal,
		Metric:      &validMeasurementUnit.Metric,
		Imperial:    &validMeasurementUnit.Imperial,
		PluralName:  &validMeasurementUnit.PluralName,
	}
}

// BuildFakeValidMeasurementUnitCreationRequestInput builds a faked ValidMeasurementUnitCreationRequestInput.
func BuildFakeValidMeasurementUnitCreationRequestInput() *types.ValidMeasurementUnitCreationRequestInput {
	validMeasurementUnit := BuildFakeValidMeasurementUnit()
	return BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(validMeasurementUnit)
}

// BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit builds a faked ValidMeasurementUnitCreationRequestInput from a valid ingredient.
func BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(validMeasurementUnit *types.ValidMeasurementUnit) *types.ValidMeasurementUnitCreationRequestInput {
	return &types.ValidMeasurementUnitCreationRequestInput{
		Name:        validMeasurementUnit.Name,
		Description: validMeasurementUnit.Description,
		Volumetric:  validMeasurementUnit.Volumetric,
		IconPath:    validMeasurementUnit.IconPath,
		Universal:   validMeasurementUnit.Universal,
		Metric:      validMeasurementUnit.Metric,
		Imperial:    validMeasurementUnit.Imperial,
		PluralName:  validMeasurementUnit.PluralName,
	}
}

// BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit builds a faked ValidMeasurementUnitDatabaseCreationInput from a valid ingredient.
func BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit(validMeasurementUnit *types.ValidMeasurementUnit) *types.ValidMeasurementUnitDatabaseCreationInput {
	return &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:          validMeasurementUnit.ID,
		Name:        validMeasurementUnit.Name,
		Description: validMeasurementUnit.Description,
		Volumetric:  validMeasurementUnit.Volumetric,
		IconPath:    validMeasurementUnit.IconPath,
		Universal:   validMeasurementUnit.Universal,
		Metric:      validMeasurementUnit.Metric,
		Imperial:    validMeasurementUnit.Imperial,
		PluralName:  validMeasurementUnit.PluralName,
	}
}
