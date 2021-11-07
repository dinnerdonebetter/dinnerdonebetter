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
		steps = append(steps, BuildFakeRecipeStep())
	}

	return &types.Recipe{
		ID:                 ksuid.New().String(),
		Name:               fake.LoremIpsumSentence(exampleQuantity),
		Source:             fake.LoremIpsumSentence(exampleQuantity),
		Description:        fake.LoremIpsumSentence(exampleQuantity),
		InspiredByRecipeID: func(x string) *string { return &x }(fake.LoremIpsumSentence(exampleQuantity)),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		CreatedByUser:      ksuid.New().String(),
		Steps:              steps,
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
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
	}
}

// BuildFakeRecipeUpdateRequestInputFromRecipe builds a faked RecipeUpdateRequestInput from a recipe.
func BuildFakeRecipeUpdateRequestInputFromRecipe(recipe *types.Recipe) *types.RecipeUpdateRequestInput {
	return &types.RecipeUpdateRequestInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
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
		ID:                 recipe.ID,
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		CreatedByUser:      recipe.CreatedByUser,
		Steps:              steps,
	}
}

// BuildFakeRecipeDatabaseCreationInput builds a faked RecipeDatabaseCreationInput.
func BuildFakeRecipeDatabaseCreationInput() *types.RecipeDatabaseCreationInput {
	recipe := BuildFakeRecipe()
	return BuildFakeRecipeDatabaseCreationInputFromRecipe(recipe)
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
		Steps:              steps,
	}
}
