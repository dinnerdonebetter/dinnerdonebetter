package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
func BuildFakeRecipeStepIngredient() *models.RecipeStepIngredient {
	return &models.RecipeStepIngredient{
		ID:                  fake.Uint64(),
		IngredientID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		QuantityType:        fake.Word(),
		QuantityValue:       fake.Float32(),
		QuantityNotes:       fake.Word(),
		ProductOfRecipe:     fake.Bool(),
		IngredientNotes:     fake.Word(),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
	}
}

// BuildFakeRecipeStepIngredientList builds a faked RecipeStepIngredientList.
func BuildFakeRecipeStepIngredientList() *models.RecipeStepIngredientList {
	exampleRecipeStepIngredient1 := BuildFakeRecipeStepIngredient()
	exampleRecipeStepIngredient2 := BuildFakeRecipeStepIngredient()
	exampleRecipeStepIngredient3 := BuildFakeRecipeStepIngredient()

	return &models.RecipeStepIngredientList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
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
		IngredientID:        recipeStepIngredient.IngredientID,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipe:     recipeStepIngredient.ProductOfRecipe,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
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
		IngredientID:        recipeStepIngredient.IngredientID,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipe:     recipeStepIngredient.ProductOfRecipe,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}
