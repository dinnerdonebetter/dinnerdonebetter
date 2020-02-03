package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomIterationMediaCreationInput creates a random IterationMediaInput
func RandomIterationMediaCreationInput() *models.IterationMediaCreationInput {
	x := &models.IterationMediaCreationInput{
		Path:              fake.Word(),
		Mimetype:          fake.Word(),
		RecipeIterationID: fake.Uint64(),
		RecipeStepID:      func(x uint64) *uint64 { return &x }(fake.Uint64()),
	}

	return x
}
