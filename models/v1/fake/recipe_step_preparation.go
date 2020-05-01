package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepPreparation builds a faked recipe step preparation.
func BuildFakeRecipeStepPreparation() *models.RecipeStepPreparation {
	return &models.RecipeStepPreparation{
		ID:                  fake.Uint64(),
		ValidPreparationID:  uint64(fake.Uint32()),
		Notes:               fake.Word(),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
	}
}

// BuildFakeRecipeStepPreparationList builds a faked RecipeStepPreparationList.
func BuildFakeRecipeStepPreparationList() *models.RecipeStepPreparationList {
	exampleRecipeStepPreparation1 := BuildFakeRecipeStepPreparation()
	exampleRecipeStepPreparation2 := BuildFakeRecipeStepPreparation()
	exampleRecipeStepPreparation3 := BuildFakeRecipeStepPreparation()

	return &models.RecipeStepPreparationList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		RecipeStepPreparations: []models.RecipeStepPreparation{
			*exampleRecipeStepPreparation1,
			*exampleRecipeStepPreparation2,
			*exampleRecipeStepPreparation3,
		},
	}
}

// BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation builds a faked RecipeStepPreparationUpdateInput from a recipe step preparation.
func BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation(recipeStepPreparation *models.RecipeStepPreparation) *models.RecipeStepPreparationUpdateInput {
	return &models.RecipeStepPreparationUpdateInput{
		ValidPreparationID:  recipeStepPreparation.ValidPreparationID,
		Notes:               recipeStepPreparation.Notes,
		BelongsToRecipeStep: recipeStepPreparation.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepPreparationCreationInput builds a faked RecipeStepPreparationCreationInput.
func BuildFakeRecipeStepPreparationCreationInput() *models.RecipeStepPreparationCreationInput {
	recipeStepPreparation := BuildFakeRecipeStepPreparation()
	return BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(recipeStepPreparation)
}

// BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation builds a faked RecipeStepPreparationCreationInput from a recipe step preparation.
func BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(recipeStepPreparation *models.RecipeStepPreparation) *models.RecipeStepPreparationCreationInput {
	return &models.RecipeStepPreparationCreationInput{
		ValidPreparationID:  recipeStepPreparation.ValidPreparationID,
		Notes:               recipeStepPreparation.Notes,
		BelongsToRecipeStep: recipeStepPreparation.BelongsToRecipeStep,
	}
}
