package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *models.Recipe {
	return &models.Recipe{
		ID:                 fake.Uint64(),
		Name:               fake.Word(),
		Source:             fake.Word(),
		Description:        fake.Word(),
		InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToUser:      fake.Uint64(),
	}
}

// BuildFakeRecipeList builds a faked RecipeList.
func BuildFakeRecipeList() *models.RecipeList {
	exampleRecipe1 := BuildFakeRecipe()
	exampleRecipe2 := BuildFakeRecipe()
	exampleRecipe3 := BuildFakeRecipe()

	return &models.RecipeList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
		},
		Recipes: []models.Recipe{
			*exampleRecipe1,
			*exampleRecipe2,
			*exampleRecipe3,
		},
	}
}

// BuildFakeRecipeUpdateInputFromRecipe builds a faked RecipeUpdateInput from a recipe.
func BuildFakeRecipeUpdateInputFromRecipe(recipe *models.Recipe) *models.RecipeUpdateInput {
	return &models.RecipeUpdateInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToUser:      recipe.BelongsToUser,
	}
}

// BuildFakeRecipeCreationInput builds a faked RecipeCreationInput.
func BuildFakeRecipeCreationInput() *models.RecipeCreationInput {
	recipe := BuildFakeRecipe()
	return BuildFakeRecipeCreationInputFromRecipe(recipe)
}

// BuildFakeRecipeCreationInputFromRecipe builds a faked RecipeCreationInput from a recipe.
func BuildFakeRecipeCreationInputFromRecipe(recipe *models.Recipe) *models.RecipeCreationInput {
	return &models.RecipeCreationInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToUser:      recipe.BelongsToUser,
	}
}
