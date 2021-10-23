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
	return &types.RecipeCreationRequestInput{
		ID:                 recipe.ID,
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToHousehold: recipe.BelongsToHousehold,
	}
}

// BuildFakeRecipeDatabaseCreationInput builds a faked RecipeDatabaseCreationInput.
func BuildFakeRecipeDatabaseCreationInput() *types.RecipeDatabaseCreationInput {
	recipe := BuildFakeRecipe()
	return BuildFakeRecipeDatabaseCreationInputFromRecipe(recipe)
}

// BuildFakeRecipeDatabaseCreationInputFromRecipe builds a faked RecipeDatabaseCreationInput from a recipe.
func BuildFakeRecipeDatabaseCreationInputFromRecipe(recipe *types.Recipe) *types.RecipeDatabaseCreationInput {
	return &types.RecipeDatabaseCreationInput{
		ID:                 recipe.ID,
		Name:               recipe.Name,
		Source:             recipe.Source,
		Description:        recipe.Description,
		InspiredByRecipeID: recipe.InspiredByRecipeID,
		BelongsToHousehold: recipe.BelongsToHousehold,
	}
}
