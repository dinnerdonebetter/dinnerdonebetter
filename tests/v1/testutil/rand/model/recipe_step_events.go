package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeStepEventCreationInput creates a random RecipeStepEventInput
func RandomRecipeStepEventCreationInput() *models.RecipeStepEventCreationInput {
	x := &models.RecipeStepEventCreationInput{
		EventType:         fake.Word(),
		Done:              fake.Bool(),
		RecipeIterationID: fake.Uint64(),
		RecipeStepID:      fake.Uint64(),
	}

	return x
}
