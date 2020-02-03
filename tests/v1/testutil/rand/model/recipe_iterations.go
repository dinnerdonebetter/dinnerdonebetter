package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeIterationCreationInput creates a random RecipeIterationInput
func RandomRecipeIterationCreationInput() *models.RecipeIterationCreationInput {
	x := &models.RecipeIterationCreationInput{
		RecipeID:            fake.Uint64(),
		EndDifficultyRating: fake.Float32(),
		EndComplexityRating: fake.Float32(),
		EndTasteRating:      fake.Float32(),
		EndOverallRating:    fake.Float32(),
	}

	return x
}
