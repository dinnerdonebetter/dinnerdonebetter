package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeCreationInput creates a random RecipeInput
func RandomRecipeCreationInput() *models.RecipeCreationInput {
	x := &models.RecipeCreationInput{
		Name:               fake.Word(),
		Source:             fake.Word(),
		Description:        fake.Word(),
		InspiredByRecipeID: func(x uint64) *uint64 { return &x }(fake.Uint64()),
	}

	return x
}
