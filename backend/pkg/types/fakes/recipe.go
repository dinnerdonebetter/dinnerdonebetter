package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *types.Recipe {
	recipeID := BuildFakeID()

	var steps []*types.RecipeStep
	for i := range exampleQuantity {
		step := BuildFakeRecipeStep()
		step.Index = uint32(i)
		step.BelongsToRecipe = recipeID
		steps = append(steps, step)
	}

	prepTasks := BuildFakeRecipePrepTasksList().Data
	for i := range prepTasks {
		prepTasks[i].BelongsToRecipe = recipeID
	}

	recipeMedia := BuildFakeRecipeMediaList().Data
	for i := range recipeMedia {
		recipeMedia[i].BelongsToRecipe = &recipeID
	}

	return &types.Recipe{
		ID:                 recipeID,
		Name:               buildUniqueString(),
		Slug:               buildUniqueString(),
		Source:             buildUniqueString(),
		Description:        buildUniqueString(),
		InspiredByRecipeID: nil,
		CreatedAt:          BuildFakeTime(),
		CreatedByUser:      BuildFakeID(),
		Steps:              steps,
		PrepTasks:          prepTasks,
		SealOfApproval:     false,
		Media:              recipeMedia,
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Max: pointer.To(float32(buildFakeNumber())),
			Min: float32(buildFakeNumber()),
		},
		PortionName:         buildUniqueString(),
		PluralPortionName:   buildUniqueString(),
		EligibleForMeals:    true,
		YieldsComponentType: "main",
		SupportingRecipes:   []*types.Recipe{},
	}
}

// BuildFakeRecipesList builds a faked RecipeList.
func BuildFakeRecipesList() *filtering.QueryFilteredResult[types.Recipe] {
	var examples []*types.Recipe
	for range exampleQuantity {
		examples = append(examples, BuildFakeRecipe())
	}

	return &filtering.QueryFilteredResult[types.Recipe]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
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
	exampleCreationInput := converters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
	examplePrepTask := BuildFakeRecipePrepTask()
	examplePrepTaskInput := converters.ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(exampleRecipe, examplePrepTask)
	examplePrepTaskInput.RecipeSteps = []*types.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		{
			BelongsToRecipeStepIndex: exampleCreationInput.Steps[0].Index,
			SatisfiesRecipeStep:      false,
		},
	}
	exampleCreationInput.PrepTasks = []*types.RecipePrepTaskWithinRecipeCreationRequestInput{
		examplePrepTaskInput,
	}

	return exampleCreationInput
}
