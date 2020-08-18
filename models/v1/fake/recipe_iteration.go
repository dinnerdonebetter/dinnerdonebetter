package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeIteration builds a faked recipe iteration.
func BuildFakeRecipeIteration() *models.RecipeIteration {
	return &models.RecipeIteration{
		ID:                  fake.Uint64(),
		RecipeID:            uint64(fake.Uint32()),
		EndDifficultyRating: fake.Float32(),
		EndComplexityRating: fake.Float32(),
		EndTasteRating:      fake.Float32(),
		EndOverallRating:    fake.Float32(),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe:     fake.Uint64(),
	}
}

// BuildFakeRecipeIterationList builds a faked RecipeIterationList.
func BuildFakeRecipeIterationList() *models.RecipeIterationList {
	exampleRecipeIteration1 := BuildFakeRecipeIteration()
	exampleRecipeIteration2 := BuildFakeRecipeIteration()
	exampleRecipeIteration3 := BuildFakeRecipeIteration()

	return &models.RecipeIterationList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
		},
		RecipeIterations: []models.RecipeIteration{
			*exampleRecipeIteration1,
			*exampleRecipeIteration2,
			*exampleRecipeIteration3,
		},
	}
}

// BuildFakeRecipeIterationUpdateInputFromRecipeIteration builds a faked RecipeIterationUpdateInput from a recipe iteration.
func BuildFakeRecipeIterationUpdateInputFromRecipeIteration(recipeIteration *models.RecipeIteration) *models.RecipeIterationUpdateInput {
	return &models.RecipeIterationUpdateInput{
		RecipeID:            recipeIteration.RecipeID,
		EndDifficultyRating: recipeIteration.EndDifficultyRating,
		EndComplexityRating: recipeIteration.EndComplexityRating,
		EndTasteRating:      recipeIteration.EndTasteRating,
		EndOverallRating:    recipeIteration.EndOverallRating,
		BelongsToRecipe:     recipeIteration.BelongsToRecipe,
	}
}

// BuildFakeRecipeIterationCreationInput builds a faked RecipeIterationCreationInput.
func BuildFakeRecipeIterationCreationInput() *models.RecipeIterationCreationInput {
	recipeIteration := BuildFakeRecipeIteration()
	return BuildFakeRecipeIterationCreationInputFromRecipeIteration(recipeIteration)
}

// BuildFakeRecipeIterationCreationInputFromRecipeIteration builds a faked RecipeIterationCreationInput from a recipe iteration.
func BuildFakeRecipeIterationCreationInputFromRecipeIteration(recipeIteration *models.RecipeIteration) *models.RecipeIterationCreationInput {
	return &models.RecipeIterationCreationInput{
		RecipeID:            recipeIteration.RecipeID,
		EndDifficultyRating: recipeIteration.EndDifficultyRating,
		EndComplexityRating: recipeIteration.EndComplexityRating,
		EndTasteRating:      recipeIteration.EndTasteRating,
		EndOverallRating:    recipeIteration.EndOverallRating,
		BelongsToRecipe:     recipeIteration.BelongsToRecipe,
	}
}
