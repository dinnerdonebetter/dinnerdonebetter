package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeValidIngredientPreparation builds a faked valid ingredient preparation.
func BuildFakeValidIngredientPreparation() *types.ValidIngredientPreparation {
	return &types.ValidIngredientPreparation{
		ID:                 ksuid.New().String(),
		Notes:              buildUniqueString(),
		ValidPreparationID: buildUniqueString(),
		ValidIngredientID:  buildUniqueString(),
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

// BuildFakeValidIngredientPreparationUpdateRequestInput builds a faked ValidIngredientPreparationUpdateRequestInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationUpdateRequestInput() *types.ValidIngredientPreparationUpdateRequestInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return &types.ValidIngredientPreparationUpdateRequestInput{
		Notes:              &validIngredientPreparation.Notes,
		ValidPreparationID: &validIngredientPreparation.ValidPreparationID,
		ValidIngredientID:  &validIngredientPreparation.ValidIngredientID,
	}
}

// BuildFakeValidIngredientPreparationUpdateRequestInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationUpdateRequestInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationUpdateRequestInputFromValidIngredientPreparation(validIngredientPreparation *types.ValidIngredientPreparation) *types.ValidIngredientPreparationUpdateRequestInput {
	return &types.ValidIngredientPreparationUpdateRequestInput{
		Notes:              &validIngredientPreparation.Notes,
		ValidPreparationID: &validIngredientPreparation.ValidPreparationID,
		ValidIngredientID:  &validIngredientPreparation.ValidIngredientID,
	}
}

// BuildFakeValidIngredientPreparationCreationRequestInput builds a faked ValidIngredientPreparationCreationRequestInput.
func BuildFakeValidIngredientPreparationCreationRequestInput() *types.ValidIngredientPreparationCreationRequestInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(validIngredientPreparation)
}

// BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationCreationRequestInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationCreationRequestInputFromValidIngredientPreparation(validIngredientPreparation *types.ValidIngredientPreparation) *types.ValidIngredientPreparationCreationRequestInput {
	return &types.ValidIngredientPreparationCreationRequestInput{
		ID:                 validIngredientPreparation.ID,
		Notes:              validIngredientPreparation.Notes,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
	}
}

// BuildFakeValidIngredientPreparationDatabaseCreationInput builds a faked ValidIngredientPreparationDatabaseCreationInput.
func BuildFakeValidIngredientPreparationDatabaseCreationInput() *types.ValidIngredientPreparationDatabaseCreationInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return BuildFakeValidIngredientPreparationDatabaseCreationInputFromValidIngredientPreparation(validIngredientPreparation)
}

// BuildFakeValidIngredientPreparationDatabaseCreationInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationDatabaseCreationInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationDatabaseCreationInputFromValidIngredientPreparation(validIngredientPreparation *types.ValidIngredientPreparation) *types.ValidIngredientPreparationDatabaseCreationInput {
	return &types.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 validIngredientPreparation.ID,
		Notes:              validIngredientPreparation.Notes,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
	}
}
