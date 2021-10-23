package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeRecipe builds a faked recipe.
func BuildFakeRecipe() *types.Recipe {
	return &types.Recipe{
		ID:                 ksuid.New().String(),
		Name:               fake.LoremIpsumSentence(exampleQuantity),
		Source:             fake.LoremIpsumSentence(exampleQuantity),
		Description:        fake.LoremIpsumSentence(exampleQuantity),
		InspiredByRecipeID: func(x string) *string { return &x }(fake.LoremIpsumSentence(exampleQuantity)),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToHousehold: fake.UUID(),
		Steps:              BuildFakeRecipeStepList().RecipeSteps,
	}
}

// BuildFakeFullRecipe builds a faked recipe.
func BuildFakeFullRecipe() *types.FullRecipe {
	var steps []*types.FullRecipeStep
	for i := 0; i < exampleQuantity; i++ {
		steps = append(steps, BuildFakeFullRecipeStep())
	}

	return &types.FullRecipe{
		ID:                 ksuid.New().String(),
		Name:               fake.Word(),
		Source:             fake.Word(),
		Description:        fake.Word(),
		InspiredByRecipeID: func(x string) *string { return &x }(fake.LoremIpsumSentence(exampleQuantity)),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToHousehold: ksuid.New().String(),
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
		BelongsToHousehold: recipe.BelongsToHousehold,
	}
}

// BuildFakeRecipeUpdateRequestInputFromRecipe builds a faked RecipeUpdateRequestInput from a recipe.
func BuildFakeRecipeUpdateRequestInputFromRecipe(recipe *types.Recipe) *types.RecipeUpdateRequestInput {
	return &types.RecipeUpdateRequestInput{
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToHousehold: recipe.BelongsToHousehold,
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
		BelongsToHousehold: recipe.BelongsToHousehold,
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
		BelongsToHousehold: recipe.BelongsToHousehold,
		Steps:              steps,
	}
}
