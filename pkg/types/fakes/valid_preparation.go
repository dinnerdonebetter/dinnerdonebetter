package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidPreparation builds a faked valid preparation.
func BuildFakeValidPreparation() *types.ValidPreparation {
	return &types.ValidPreparation{
		ID:          uint64(fake.Uint32()),
		ExternalID:  fake.UUID(),
		Name:        fake.Word(),
		Description: fake.Word(),
		IconPath:    fake.Word(),
		CreatedOn:   uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidPreparationList builds a faked ValidPreparationList.
func BuildFakeValidPreparationList() *types.ValidPreparationList {
	var examples []*types.ValidPreparation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidPreparation())
	}

	return &types.ValidPreparationList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidPreparations: examples,
	}
}

// BuildFakeValidPreparationUpdateInput builds a faked ValidPreparationUpdateInput from a valid preparation.
func BuildFakeValidPreparationUpdateInput() *types.ValidPreparationUpdateInput {
	validPreparation := BuildFakeValidPreparation()
	return &types.ValidPreparationUpdateInput{
		Name:        validPreparation.Name,
		Description: validPreparation.Description,
		IconPath:    validPreparation.IconPath,
	}
}

// BuildFakeValidPreparationUpdateInputFromValidPreparation builds a faked ValidPreparationUpdateInput from a valid preparation.
func BuildFakeValidPreparationUpdateInputFromValidPreparation(validPreparation *types.ValidPreparation) *types.ValidPreparationUpdateInput {
	return &types.ValidPreparationUpdateInput{
		Name:        validPreparation.Name,
		Description: validPreparation.Description,
		IconPath:    validPreparation.IconPath,
	}
}

// BuildFakeValidPreparationCreationInput builds a faked ValidPreparationCreationInput.
func BuildFakeValidPreparationCreationInput() *types.ValidPreparationCreationInput {
	validPreparation := BuildFakeValidPreparation()
	return BuildFakeValidPreparationCreationInputFromValidPreparation(validPreparation)
}

// BuildFakeValidPreparationCreationInputFromValidPreparation builds a faked ValidPreparationCreationInput from a valid preparation.
func BuildFakeValidPreparationCreationInputFromValidPreparation(validPreparation *types.ValidPreparation) *types.ValidPreparationCreationInput {
	return &types.ValidPreparationCreationInput{
		Name:        validPreparation.Name,
		Description: validPreparation.Description,
		IconPath:    validPreparation.IconPath,
	}
}
