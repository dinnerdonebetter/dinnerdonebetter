package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeMeal builds a faked recipe.
func BuildFakeMeal() *types.Meal {
	recipes := []*types.Recipe{}
	for i := 0; i < exampleQuantity; i++ {
		recipes = append(recipes, BuildFakeRecipe())
	}

	return &types.Meal{
		ID:            ksuid.New().String(),
		Name:          fake.LoremIpsumSentence(exampleQuantity),
		Description:   fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:     uint64(uint32(fake.Date().Unix())),
		CreatedByUser: ksuid.New().String(),
		Recipes:       recipes,
	}
}

// BuildFakeMealList builds a faked MealList.
func BuildFakeMealList() *types.MealList {
	var examples []*types.Meal
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMeal())
	}

	return &types.MealList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Meals: examples,
	}
}

// BuildFakeMealUpdateRequestInput builds a faked MealUpdateRequestInput from a recipe.
func BuildFakeMealUpdateRequestInput() *types.MealUpdateRequestInput {
	recipe := BuildFakeMeal()

	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealUpdateRequestInput{
		Name:          recipe.Name,
		Description:   recipe.Description,
		CreatedByUser: recipe.CreatedByUser,
		Recipes:       recipeIDs,
	}
}

// BuildFakeMealUpdateRequestInputFromMeal builds a faked MealUpdateRequestInput from a recipe.
func BuildFakeMealUpdateRequestInputFromMeal(recipe *types.Meal) *types.MealUpdateRequestInput {
	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealUpdateRequestInput{
		Name:          recipe.Name,
		Description:   recipe.Description,
		CreatedByUser: recipe.CreatedByUser,
		Recipes:       recipeIDs,
	}
}

// BuildFakeMealCreationRequestInput builds a faked MealCreationRequestInput.
func BuildFakeMealCreationRequestInput() *types.MealCreationRequestInput {
	recipe := BuildFakeMeal()
	return BuildFakeMealCreationRequestInputFromMeal(recipe)
}

// BuildFakeMealCreationRequestInputFromMeal builds a faked MealCreationRequestInput from a recipe.
func BuildFakeMealCreationRequestInputFromMeal(recipe *types.Meal) *types.MealCreationRequestInput {
	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealCreationRequestInput{
		ID:            recipe.ID,
		Name:          recipe.Name,
		Description:   recipe.Description,
		CreatedByUser: recipe.CreatedByUser,
		Recipes:       recipeIDs,
	}
}

// BuildFakeMealDatabaseCreationInput builds a faked MealDatabaseCreationInput.
func BuildFakeMealDatabaseCreationInput() *types.MealDatabaseCreationInput {
	recipe := BuildFakeMeal()
	return BuildFakeMealDatabaseCreationInputFromMeal(recipe)
}

// BuildFakeMealDatabaseCreationInputFromMeal builds a faked MealDatabaseCreationInput from a recipe.
func BuildFakeMealDatabaseCreationInputFromMeal(recipe *types.Meal) *types.MealDatabaseCreationInput {
	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealDatabaseCreationInput{
		ID:            recipe.ID,
		Name:          recipe.Name,
		Description:   recipe.Description,
		CreatedByUser: recipe.CreatedByUser,
		Recipes:       recipeIDs,
	}
}
