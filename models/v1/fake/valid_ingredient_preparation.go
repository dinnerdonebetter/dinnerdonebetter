package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidIngredientPreparation builds a faked valid ingredient preparation.
func BuildFakeValidIngredientPreparation() *models.ValidIngredientPreparation {
	return &models.ValidIngredientPreparation{
		ID:                 fake.Uint64(),
		Notes:              fake.Word(),
		ValidPreparationID: uint64(fake.Uint32()),
		ValidIngredientID:  uint64(fake.Uint32()),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidIngredientPreparationList builds a faked ValidIngredientPreparationList.
func BuildFakeValidIngredientPreparationList() *models.ValidIngredientPreparationList {
	exampleValidIngredientPreparation1 := BuildFakeValidIngredientPreparation()
	exampleValidIngredientPreparation2 := BuildFakeValidIngredientPreparation()
	exampleValidIngredientPreparation3 := BuildFakeValidIngredientPreparation()

	return &models.ValidIngredientPreparationList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		ValidIngredientPreparations: []models.ValidIngredientPreparation{
			*exampleValidIngredientPreparation1,
			*exampleValidIngredientPreparation2,
			*exampleValidIngredientPreparation3,
		},
	}
}

// BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationUpdateInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(validIngredientPreparation *models.ValidIngredientPreparation) *models.ValidIngredientPreparationUpdateInput {
	return &models.ValidIngredientPreparationUpdateInput{
		Notes:              validIngredientPreparation.Notes,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
	}
}

// BuildFakeValidIngredientPreparationCreationInput builds a faked ValidIngredientPreparationCreationInput.
func BuildFakeValidIngredientPreparationCreationInput() *models.ValidIngredientPreparationCreationInput {
	validIngredientPreparation := BuildFakeValidIngredientPreparation()
	return BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(validIngredientPreparation)
}

// BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation builds a faked ValidIngredientPreparationCreationInput from a valid ingredient preparation.
func BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(validIngredientPreparation *models.ValidIngredientPreparation) *models.ValidIngredientPreparationCreationInput {
	return &models.ValidIngredientPreparationCreationInput{
		Notes:              validIngredientPreparation.Notes,
		ValidPreparationID: validIngredientPreparation.ValidPreparationID,
		ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
	}
}
