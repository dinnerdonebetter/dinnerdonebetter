package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeStepCreationInput creates a random RecipeStepInput
func RandomRecipeStepCreationInput() *models.RecipeStepCreationInput {
	x := &models.RecipeStepCreationInput{
		Index:                     uint(fake.Uint32()),
		PreparationID:             fake.Uint64(),
		PrerequisiteStep:          fake.Uint64(),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
		Notes:                     fake.Word(),
		RecipeID:                  fake.Uint64(),
	}

	return x
}
