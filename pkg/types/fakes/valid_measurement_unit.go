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
		Name:        fake.LoremIpsumSentence(exampleQuantity),
		Description: fake.LoremIpsumSentence(exampleQuantity),
		IconPath:    fake.LoremIpsumSentence(exampleQuantity),
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
	validIngredient := BuildFakeValidMeasurementUnit()
	return &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        &validIngredient.Name,
		Description: &validIngredient.Description,
		IconPath:    &validIngredient.IconPath,
	}
}

// BuildFakeValidMeasurementUnitUpdateRequestInputFromValidMeasurementUnit builds a faked ValidMeasurementUnitUpdateRequestInput from a valid ingredient.
func BuildFakeValidMeasurementUnitUpdateRequestInputFromValidMeasurementUnit(validIngredient *types.ValidMeasurementUnit) *types.ValidMeasurementUnitUpdateRequestInput {
	return &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        &validIngredient.Name,
		Description: &validIngredient.Description,
		IconPath:    &validIngredient.IconPath,
	}
}

// BuildFakeValidMeasurementUnitCreationRequestInput builds a faked ValidMeasurementUnitCreationRequestInput.
func BuildFakeValidMeasurementUnitCreationRequestInput() *types.ValidMeasurementUnitCreationRequestInput {
	validIngredient := BuildFakeValidMeasurementUnit()
	return BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(validIngredient)
}

// BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit builds a faked ValidMeasurementUnitCreationRequestInput from a valid ingredient.
func BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(validIngredient *types.ValidMeasurementUnit) *types.ValidMeasurementUnitCreationRequestInput {
	return &types.ValidMeasurementUnitCreationRequestInput{
		ID:          validIngredient.ID,
		Name:        validIngredient.Name,
		Description: validIngredient.Description,
		IconPath:    validIngredient.IconPath,
	}
}

// BuildFakeValidMeasurementUnitDatabaseCreationInput builds a faked ValidMeasurementUnitDatabaseCreationInput.
func BuildFakeValidMeasurementUnitDatabaseCreationInput() *types.ValidMeasurementUnitDatabaseCreationInput {
	validIngredient := BuildFakeValidMeasurementUnit()
	return BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit(validIngredient)
}

// BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit builds a faked ValidMeasurementUnitDatabaseCreationInput from a valid ingredient.
func BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit(validIngredient *types.ValidMeasurementUnit) *types.ValidMeasurementUnitDatabaseCreationInput {
	return &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:          validIngredient.ID,
		Name:        validIngredient.Name,
		Description: validIngredient.Description,
		IconPath:    validIngredient.IconPath,
	}
}
