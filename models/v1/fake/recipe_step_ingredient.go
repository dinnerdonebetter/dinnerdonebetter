package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
func BuildFakeRecipeStepIngredient() *models.RecipeStepIngredient {
	return &models.RecipeStepIngredient{
		ID:                    fake.Uint64(),
		ValidIngredientID:     uint64(fake.Uint32()),
		IngredientNotes:       fake.Word(),
		QuantityType:          fake.Word(),
		QuantityValue:         fake.Float32(),
		QuantityNotes:         fake.Word(),
		ProductOfRecipeStepID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		CreatedOn:             uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep:   fake.Uint64(),
	}
}

// BuildFakeRecipeStepIngredientList builds a faked RecipeStepIngredientList.
func BuildFakeRecipeStepIngredientList() *models.RecipeStepIngredientList {
	exampleRecipeStepIngredient1 := BuildFakeRecipeStepIngredient()
	exampleRecipeStepIngredient2 := BuildFakeRecipeStepIngredient()
	exampleRecipeStepIngredient3 := BuildFakeRecipeStepIngredient()

	return &models.RecipeStepIngredientList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		RecipeStepIngredients: []models.RecipeStepIngredient{
			*exampleRecipeStepIngredient1,
			*exampleRecipeStepIngredient2,
			*exampleRecipeStepIngredient3,
		},
	}
}

// BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient builds a faked RecipeStepIngredientUpdateInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(recipeStepIngredient *models.RecipeStepIngredient) *models.RecipeStepIngredientUpdateInput {
	return &models.RecipeStepIngredientUpdateInput{
		ValidIngredientID:     recipeStepIngredient.ValidIngredientID,
		IngredientNotes:       recipeStepIngredient.IngredientNotes,
		QuantityType:          recipeStepIngredient.QuantityType,
		QuantityValue:         recipeStepIngredient.QuantityValue,
		QuantityNotes:         recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStepID: recipeStepIngredient.ProductOfRecipeStepID,
		BelongsToRecipeStep:   recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientCreationInput builds a faked RecipeStepIngredientCreationInput.
func BuildFakeRecipeStepIngredientCreationInput() *models.RecipeStepIngredientCreationInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(recipeStepIngredient)
}

// BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient builds a faked RecipeStepIngredientCreationInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(recipeStepIngredient *models.RecipeStepIngredient) *models.RecipeStepIngredientCreationInput {
	return &models.RecipeStepIngredientCreationInput{
		ValidIngredientID:     recipeStepIngredient.ValidIngredientID,
		IngredientNotes:       recipeStepIngredient.IngredientNotes,
		QuantityType:          recipeStepIngredient.QuantityType,
		QuantityValue:         recipeStepIngredient.QuantityValue,
		QuantityNotes:         recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStepID: recipeStepIngredient.ProductOfRecipeStepID,
		BelongsToRecipeStep:   recipeStepIngredient.BelongsToRecipeStep,
	}
}
