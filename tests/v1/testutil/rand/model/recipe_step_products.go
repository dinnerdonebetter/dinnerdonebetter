package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeStepProductCreationInput creates a random RecipeStepProductInput
func RandomRecipeStepProductCreationInput() *models.RecipeStepProductCreationInput {
	x := &models.RecipeStepProductCreationInput{
		Name:         fake.Word(),
		RecipeStepID: fake.Uint64(),
	}

	return x
}
