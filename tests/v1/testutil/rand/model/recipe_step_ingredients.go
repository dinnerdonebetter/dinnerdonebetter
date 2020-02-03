package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRecipeStepIngredientCreationInput creates a random RecipeStepIngredientInput
func RandomRecipeStepIngredientCreationInput() *models.RecipeStepIngredientCreationInput {
	x := &models.RecipeStepIngredientCreationInput{
		IngredientID:    func(x uint64) *uint64 { return &x }(fake.Uint64()),
		QuantityType:    fake.Word(),
		QuantityValue:   fake.Float32(),
		QuantityNotes:   fake.Word(),
		ProductOfRecipe: fake.Bool(),
		IngredientNotes: fake.Word(),
		RecipeStepID:    fake.Uint64(),
	}

	return x
}
