package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomPreparationCreationInput creates a random PreparationInput
func RandomPreparationCreationInput() *models.PreparationCreationInput {
	x := &models.PreparationCreationInput{
		Name:           fake.Word(),
		Variant:        fake.Word(),
		Description:    fake.Word(),
		AllergyWarning: fake.Word(),
		Icon:           fake.Word(),
	}

	return x
}
