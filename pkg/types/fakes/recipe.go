package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *types.Recipe {
	recipeID := BuildFakeID()

	var steps []*types.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		step := BuildFakeRecipeStep()
		step.Index = uint32(i)
		step.BelongsToRecipe = recipeID
		steps = append(steps, step)
	}

	prepTasks := BuildFakeRecipePrepTaskList().RecipePrepTasks
	for i := range prepTasks {
		prepTasks[i].BelongsToRecipe = recipeID
	}

	recipeMedia := BuildFakeRecipeMediaList().RecipeMedia
	for i := range recipeMedia {
		recipeMedia[i].BelongsToRecipe = &recipeID
	}

	return &types.Recipe{
		ID:                 recipeID,
		Name:               buildUniqueString(),
		Source:             buildUniqueString(),
		Description:        buildUniqueString(),
		InspiredByRecipeID: func(x string) *string { return &x }(buildUniqueString()),
		CreatedAt:          BuildFakeTime(),
		CreatedByUser:      BuildFakeID(),
		Steps:              steps,
		PrepTasks:          prepTasks,
		SealOfApproval:     false,
		Media:              recipeMedia,
		YieldsPortions:     fake.Uint8(),
	}
}

// BuildFakeRecipeList builds a faked RecipeList.
func BuildFakeRecipeList() *types.RecipeList {
	var examples []*types.Recipe
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipe())
	}

	return &types.RecipeList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Recipes: examples,
	}
}

// BuildFakeRecipeUpdateRequestInput builds a faked RecipeUpdateRequestInput from a recipe.
func BuildFakeRecipeUpdateRequestInput() *types.RecipeUpdateRequestInput {
	recipe := BuildFakeRecipe()
	return converters.ConvertRecipeToRecipeUpdateRequestInput(recipe)
}

// BuildFakeRecipeCreationRequestInput builds a faked RecipeCreationRequestInput.
func BuildFakeRecipeCreationRequestInput() *types.RecipeCreationRequestInput {
	exampleRecipe := BuildFakeRecipe()
	exampleCreationInput := converters.ConvertRecipeToRecipeCreationRequestInputFromRecipe(exampleRecipe)
	exampleCreationInput.CreatedByUser = ""
	examplePrepTask := BuildFakeRecipePrepTask()
	examplePrepTaskInput := converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask)
	examplePrepTaskInput.TaskSteps = []*types.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		{
			BelongsToRecipeStepIndex: exampleCreationInput.Steps[0].Index,
			BelongsToRecipePrepTask:  examplePrepTask.ID,
			SatisfiesRecipeStep:      false,
		},
	}
	exampleCreationInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
		examplePrepTaskInput,
	}

	return exampleCreationInput
}
