package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeIterationStep builds a faked recipe iteration step.
func BuildFakeRecipeIterationStep() *models.RecipeIterationStep {
	return &models.RecipeIterationStep{
		ID:              fake.Uint64(),
		StartedOn:       func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		EndedOn:         func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		State:           fake.Word(),
		CreatedOn:       uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe: fake.Uint64(),
	}
}

// BuildFakeRecipeIterationStepList builds a faked RecipeIterationStepList.
func BuildFakeRecipeIterationStepList() *models.RecipeIterationStepList {
	exampleRecipeIterationStep1 := BuildFakeRecipeIterationStep()
	exampleRecipeIterationStep2 := BuildFakeRecipeIterationStep()
	exampleRecipeIterationStep3 := BuildFakeRecipeIterationStep()

	return &models.RecipeIterationStepList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		RecipeIterationSteps: []models.RecipeIterationStep{
			*exampleRecipeIterationStep1,
			*exampleRecipeIterationStep2,
			*exampleRecipeIterationStep3,
		},
	}
}

// BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep builds a faked RecipeIterationStepUpdateInput from a recipe iteration step.
func BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep(recipeIterationStep *models.RecipeIterationStep) *models.RecipeIterationStepUpdateInput {
	return &models.RecipeIterationStepUpdateInput{
		StartedOn:       recipeIterationStep.StartedOn,
		EndedOn:         recipeIterationStep.EndedOn,
		State:           recipeIterationStep.State,
		BelongsToRecipe: recipeIterationStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeIterationStepCreationInput builds a faked RecipeIterationStepCreationInput.
func BuildFakeRecipeIterationStepCreationInput() *models.RecipeIterationStepCreationInput {
	recipeIterationStep := BuildFakeRecipeIterationStep()
	return BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(recipeIterationStep)
}

// BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep builds a faked RecipeIterationStepCreationInput from a recipe iteration step.
func BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(recipeIterationStep *models.RecipeIterationStep) *models.RecipeIterationStepCreationInput {
	return &models.RecipeIterationStepCreationInput{
		StartedOn:       recipeIterationStep.StartedOn,
		EndedOn:         recipeIterationStep.EndedOn,
		State:           recipeIterationStep.State,
		BelongsToRecipe: recipeIterationStep.BelongsToRecipe,
	}
}
