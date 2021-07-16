package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidIngredientPreparation builds a faked valid ingredient preparation.
func BuildFakeValidIngredientPreparation() *types.ValidIngredientPreparation {
	return &types.ValidIngredientPreparation{
		ID:                 uint64(fake.Uint32()),
		ExternalID:         fake.UUID(),
		Notes:              fake.Word(),
		ValidIngredientID:  uint64(fake.Uint32()),
		ValidPreparationID: uint64(fake.Uint32()),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidIngredientPreparationList builds a faked ValidIngredientPreparationList.
func BuildFakeValidIngredientPreparationList() *types.ValidIngredientPreparationList {
	var examples []*types.ValidIngredientPreparation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientPreparation())
	}

	return &types.ValidIngredientPreparationList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidIngredientPreparations: examples,
	}
}

// BuildFakeValidIngredientPreparationUpdateInput builds a faked ValidIngredientPreparationUpdateInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationUpdateInput() *types.ValidIngredientPreparationUpdateInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return &types.ValidIngredientPreparationUpdateInput{
		Notes:              validIngredientPreparation.Notes,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
	}
}

// BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationUpdateInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(validIngredientPreparation *types.ValidIngredientPreparation) *types.ValidIngredientPreparationUpdateInput {
	return &types.ValidIngredientPreparationUpdateInput{
		Notes:              validIngredientPreparation.Notes,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
	}
}

// BuildFakeValidIngredientPreparationCreationInput builds a faked ValidIngredientPreparationCreationInput.
func BuildFakeValidIngredientPreparationCreationInput() *types.ValidIngredientPreparationCreationInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(validIngredientPreparation)
}

// BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationCreationInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(validIngredientPreparation *types.ValidIngredientPreparation) *types.ValidIngredientPreparationCreationInput {
	return &types.ValidIngredientPreparationCreationInput{
		Notes:              validIngredientPreparation.Notes,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
	}
}
