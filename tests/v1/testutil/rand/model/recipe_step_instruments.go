package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeStepInstrumentCreationInput creates a random RecipeStepInstrumentInput
func RandomRecipeStepInstrumentCreationInput() *models.RecipeStepInstrumentCreationInput {
	x := &models.RecipeStepInstrumentCreationInput{
		InstrumentID: func(x uint64) *uint64 { return &x }(fake.Uint64()),
		RecipeStepID: fake.Uint64(),
		Notes:        fake.Word(),
	}

	return x
}
