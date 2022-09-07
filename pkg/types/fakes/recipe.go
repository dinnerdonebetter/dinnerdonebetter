package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *types.Recipe {
	var steps []*types.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		step := BuildFakeRecipeStep()
		step.Index = uint32(i)
		steps = append(steps, step)
	}

	return &types.Recipe{
		ID:                 ksuid.New().String(),
		Name:               buildUniqueString(),
		Source:             buildUniqueString(),
		Description:        buildUniqueString(),
		InspiredByRecipeID: func(x string) *string { return &x }(buildUniqueString()),
		CreatedAt:          fake.Date(),
		CreatedByUser:      ksuid.New().String(),
		Steps:              steps,
		SealOfApproval:     false,
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
	return &types.RecipeUpdateRequestInput{
		Name:               &recipe.Name,
		Source:             &recipe.Source,
		Description:        &recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      &recipe.CreatedByUser,
		SealOfApproval:     &recipe.SealOfApproval,
		YieldsPortions:     &recipe.YieldsPortions,
	}
}

// BuildFakeRecipeUpdateRequestInputFromRecipe builds a faked RecipeUpdateRequestInput from a recipe.
func BuildFakeRecipeUpdateRequestInputFromRecipe(recipe *types.Recipe) *types.RecipeUpdateRequestInput {
	return &types.RecipeUpdateRequestInput{
		Name:               &recipe.Name,
		Source:             &recipe.Source,
		Description:        &recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      &recipe.CreatedByUser,
		SealOfApproval:     &recipe.SealOfApproval,
		YieldsPortions:     &recipe.YieldsPortions,
	}
}

// BuildFakeRecipeCreationRequestInput builds a faked RecipeCreationRequestInput.
func BuildFakeRecipeCreationRequestInput() *types.RecipeCreationRequestInput {
	recipe := BuildFakeRecipe()
	return BuildFakeRecipeCreationRequestInputFromRecipe(recipe)
}

// BuildFakeRecipeCreationRequestInputFromRecipe builds a faked RecipeCreationRequestInput from a recipe.
func BuildFakeRecipeCreationRequestInputFromRecipe(recipe *types.Recipe) *types.RecipeCreationRequestInput {
	steps := []*types.RecipeStepCreationRequestInput{}
	for _, step := range recipe.Steps {
		steps = append(steps, BuildFakeRecipeStepCreationRequestInputFromRecipeStep(step))
	}

	return &types.RecipeCreationRequestInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
		SealOfApproval:     recipe.SealOfApproval,
		YieldsPortions:     recipe.YieldsPortions,
		Steps:              steps,
	}
}

// BuildFakeRecipeDatabaseCreationInputFromRecipe builds a faked RecipeDatabaseCreationInput from a recipe.
func BuildFakeRecipeDatabaseCreationInputFromRecipe(recipe *types.Recipe) *types.RecipeDatabaseCreationInput {
	steps := []*types.RecipeStepDatabaseCreationInput{}
	for _, step := range recipe.Steps {
		steps = append(steps, BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(step))
	}

	return &types.RecipeDatabaseCreationInput{
		ID:                 recipe.ID,
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
		SealOfApproval:     recipe.SealOfApproval,
		YieldsPortions:     recipe.YieldsPortions,
		Steps:              steps,
	}
}
