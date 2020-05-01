package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeTag builds a faked recipe tag.
func BuildFakeRecipeTag() *models.RecipeTag {
	return &models.RecipeTag{
		ID:              fake.Uint64(),
		Name:            fake.Word(),
		CreatedOn:       uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe: fake.Uint64(),
	}
}

// BuildFakeRecipeTagList builds a faked RecipeTagList.
func BuildFakeRecipeTagList() *models.RecipeTagList {
	exampleRecipeTag1 := BuildFakeRecipeTag()
	exampleRecipeTag2 := BuildFakeRecipeTag()
	exampleRecipeTag3 := BuildFakeRecipeTag()

	return &models.RecipeTagList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		RecipeTags: []models.RecipeTag{
			*exampleRecipeTag1,
			*exampleRecipeTag2,
			*exampleRecipeTag3,
		},
	}
}

// BuildFakeRecipeTagUpdateInputFromRecipeTag builds a faked RecipeTagUpdateInput from a recipe tag.
func BuildFakeRecipeTagUpdateInputFromRecipeTag(recipeTag *models.RecipeTag) *models.RecipeTagUpdateInput {
	return &models.RecipeTagUpdateInput{
		Name:            recipeTag.Name,
		BelongsToRecipe: recipeTag.BelongsToRecipe,
	}
}

// BuildFakeRecipeTagCreationInput builds a faked RecipeTagCreationInput.
func BuildFakeRecipeTagCreationInput() *models.RecipeTagCreationInput {
	recipeTag := BuildFakeRecipeTag()
	return BuildFakeRecipeTagCreationInputFromRecipeTag(recipeTag)
}

// BuildFakeRecipeTagCreationInputFromRecipeTag builds a faked RecipeTagCreationInput from a recipe tag.
func BuildFakeRecipeTagCreationInputFromRecipeTag(recipeTag *models.RecipeTag) *models.RecipeTagCreationInput {
	return &models.RecipeTagCreationInput{
		Name:            recipeTag.Name,
		BelongsToRecipe: recipeTag.BelongsToRecipe,
	}
}
