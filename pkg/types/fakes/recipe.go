package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *types.Recipe {
	return &types.Recipe{
		ID:          uint64(fake.Uint32()),
		ExternalID:  fake.UUID(),
		Name:        fake.Word(),
		Source:      fake.Word(),
		Description: fake.Word(),
		//DisplayImageURL:    fake.ImageURL(123, 123),
		InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToHousehold: fake.Uint64(),
		Steps:              BuildFakeRecipeStepList().RecipeSteps,
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

// BuildFakeRecipeUpdateInput builds a faked RecipeUpdateInput from a recipe.
func BuildFakeRecipeUpdateInput() *types.RecipeUpdateInput {
	recipe := BuildFakeRecipe()
	return &types.RecipeUpdateInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		DisplayImageURL:    recipe.DisplayImageURL,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToHousehold: recipe.BelongsToHousehold,
	}
}

// BuildFakeRecipeUpdateInputFromRecipe builds a faked RecipeUpdateInput from a recipe.
func BuildFakeRecipeUpdateInputFromRecipe(recipe *types.Recipe) *types.RecipeUpdateInput {
	return &types.RecipeUpdateInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		DisplayImageURL:    recipe.DisplayImageURL,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToHousehold: recipe.BelongsToHousehold,
	}
}

// BuildFakeRecipeCreationInput builds a faked RecipeCreationInput.
func BuildFakeRecipeCreationInput() *types.RecipeCreationInput {
	recipe := BuildFakeRecipe()
	return BuildFakeRecipeCreationInputFromRecipe(recipe)
}

// BuildFakeRecipeCreationInputFromRecipe builds a faked RecipeCreationInput from a recipe.
func BuildFakeRecipeCreationInputFromRecipe(recipe *types.Recipe) *types.RecipeCreationInput {
	steps := []*types.RecipeStepCreationInput{}
	for _, step := range recipe.Steps {
		steps = append(steps, BuildFakeRecipeStepCreationInputFromRecipeStep(step))
	}

	return &types.RecipeCreationInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		DisplayImageURL:    recipe.DisplayImageURL,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToHousehold: recipe.BelongsToHousehold,
		Steps:              steps,
	}
}
