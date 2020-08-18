package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidPreparation builds a faked valid preparation.
func BuildFakeValidPreparation() *models.ValidPreparation {
	return &models.ValidPreparation{
		ID:                         fake.Uint64(),
		Name:                       fake.Word(),
		Description:                fake.Word(),
		Icon:                       fake.Word(),
		ApplicableToAllIngredients: fake.Bool(),
		CreatedOn:                  uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidPreparationList builds a faked ValidPreparationList.
func BuildFakeValidPreparationList() *models.ValidPreparationList {
	exampleValidPreparation1 := BuildFakeValidPreparation()
	exampleValidPreparation2 := BuildFakeValidPreparation()
	exampleValidPreparation3 := BuildFakeValidPreparation()

	return &models.ValidPreparationList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		ValidPreparations: []models.ValidPreparation{
			*exampleValidPreparation1,
			*exampleValidPreparation2,
			*exampleValidPreparation3,
		},
	}
}

// BuildFakeValidPreparationUpdateInputFromValidPreparation builds a faked ValidPreparationUpdateInput from a valid preparation.
func BuildFakeValidPreparationUpdateInputFromValidPreparation(validPreparation *models.ValidPreparation) *models.ValidPreparationUpdateInput {
	return &models.ValidPreparationUpdateInput{
		Name:                       validPreparation.Name,
		Description:                validPreparation.Description,
		Icon:                       validPreparation.Icon,
		ApplicableToAllIngredients: validPreparation.ApplicableToAllIngredients,
	}
}

// BuildFakeValidPreparationCreationInput builds a faked ValidPreparationCreationInput.
func BuildFakeValidPreparationCreationInput() *models.ValidPreparationCreationInput {
	validPreparation := BuildFakeValidPreparation()
	return BuildFakeValidPreparationCreationInputFromValidPreparation(validPreparation)
}

// BuildFakeValidPreparationCreationInputFromValidPreparation builds a faked ValidPreparationCreationInput from a valid preparation.
func BuildFakeValidPreparationCreationInputFromValidPreparation(validPreparation *models.ValidPreparation) *models.ValidPreparationCreationInput {
	return &models.ValidPreparationCreationInput{
		Name:                       validPreparation.Name,
		Description:                validPreparation.Description,
		Icon:                       validPreparation.Icon,
		ApplicableToAllIngredients: validPreparation.ApplicableToAllIngredients,
	}
}
