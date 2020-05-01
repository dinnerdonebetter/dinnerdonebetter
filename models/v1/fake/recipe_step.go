package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *models.RecipeStep {
	return &models.RecipeStep{
		ID:                        fake.Uint64(),
		Index:                     uint(fake.Uint32()),
		ValidPreparationID:        uint64(fake.Uint32()),
		PrerequisiteStepID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		YieldsProductName:         fake.Word(),
		YieldsQuantity:            uint(fake.Uint32()),
		Notes:                     fake.Word(),
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe:           fake.Uint64(),
	}
}

// BuildFakeRecipeStepList builds a faked RecipeStepList.
func BuildFakeRecipeStepList() *models.RecipeStepList {
	exampleRecipeStep1 := BuildFakeRecipeStep()
	exampleRecipeStep2 := BuildFakeRecipeStep()
	exampleRecipeStep3 := BuildFakeRecipeStep()

	return &models.RecipeStepList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		RecipeSteps: []models.RecipeStep{
			*exampleRecipeStep1,
			*exampleRecipeStep2,
			*exampleRecipeStep3,
		},
	}
}

// BuildFakeRecipeStepUpdateInputFromRecipeStep builds a faked RecipeStepUpdateInput from a recipe step.
func BuildFakeRecipeStepUpdateInputFromRecipeStep(recipeStep *models.RecipeStep) *models.RecipeStepUpdateInput {
	return &models.RecipeStepUpdateInput{
		Index:                     recipeStep.Index,
		ValidPreparationID:        recipeStep.ValidPreparationID,
		PrerequisiteStepID:        recipeStep.PrerequisiteStepID,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		YieldsProductName:         recipeStep.YieldsProductName,
		YieldsQuantity:            recipeStep.YieldsQuantity,
		Notes:                     recipeStep.Notes,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepCreationInput builds a faked RecipeStepCreationInput.
func BuildFakeRecipeStepCreationInput() *models.RecipeStepCreationInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepCreationInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepCreationInputFromRecipeStep builds a faked RecipeStepCreationInput from a recipe step.
func BuildFakeRecipeStepCreationInputFromRecipeStep(recipeStep *models.RecipeStep) *models.RecipeStepCreationInput {
	return &models.RecipeStepCreationInput{
		Index:                     recipeStep.Index,
		ValidPreparationID:        recipeStep.ValidPreparationID,
		PrerequisiteStepID:        recipeStep.PrerequisiteStepID,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		YieldsProductName:         recipeStep.YieldsProductName,
		YieldsQuantity:            recipeStep.YieldsQuantity,
		Notes:                     recipeStep.Notes,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}
