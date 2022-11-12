package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeValidIngredientPreparation builds a faked valid ingredient preparation.
func BuildFakeValidIngredientPreparation() *types.ValidIngredientPreparation {
	return &types.ValidIngredientPreparation{
		ID:          BuildFakeID(),
		Notes:       buildUniqueString(),
		Preparation: *BuildFakeValidPreparation(),
		Ingredient:  *BuildFakeValidIngredient(),
		CreatedAt:   BuildFakeTime(),
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
		ValidPreparationID: &validIngredientPreparation.Preparation.ID,
		ValidIngredientID:  &validIngredientPreparation.Ingredient.ID,
	}
}

// BuildFakeValidIngredientPreparationCreationRequestInput builds a faked ValidIngredientPreparationCreationRequestInput.
func BuildFakeValidIngredientPreparationCreationRequestInput() *types.ValidIngredientPreparationCreationRequestInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(validIngredientPreparation)
}
