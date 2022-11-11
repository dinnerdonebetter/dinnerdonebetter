package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
)

// BuildFakeMeal builds a faked recipe.
func BuildFakeMeal() *types.Meal {
	recipes := []*types.Recipe{}
	for i := 0; i < exampleQuantity; i++ {
		recipes = append(recipes, BuildFakeRecipe())
	}

	return &types.Meal{
		ID:            BuildFakeID(),
		Name:          buildUniqueString(),
		Description:   buildUniqueString(),
		CreatedAt:     BuildFakeTime(),
		CreatedByUser: BuildFakeID(),
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

// BuildFakeMealUpdateRequestInputFromMeal builds a faked MealUpdateRequestInput from a recipe.
func BuildFakeMealUpdateRequestInputFromMeal(recipe *types.Meal) *types.MealUpdateRequestInput {
	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealUpdateRequestInput{
		Name:          &recipe.Name,
		Description:   &recipe.Description,
		CreatedByUser: &recipe.CreatedByUser,
		Recipes:       recipeIDs,
	}
}

// BuildFakeMealCreationRequestInput builds a faked MealCreationRequestInput.
func BuildFakeMealCreationRequestInput() *types.MealCreationRequestInput {
	recipe := BuildFakeMeal()
	return BuildFakeMealCreationRequestInputFromMeal(recipe)
}

// BuildFakeMealCreationRequestInputFromMeal builds a faked MealCreationRequestInput from a recipe.
func BuildFakeMealCreationRequestInputFromMeal(meal *types.Meal) *types.MealCreationRequestInput {
	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealCreationRequestInput{
		ID:            meal.ID,
		Name:          meal.Name,
		Description:   meal.Description,
		CreatedByUser: meal.CreatedByUser,
		Recipes:       recipeIDs,
	}
}

// ConvertMealToMealDatabaseCreationInput builds a faked MealDatabaseCreationInput from a recipe.
func ConvertMealToMealDatabaseCreationInput(meal *types.Meal) *types.MealDatabaseCreationInput {
	recipeIDs := []string{}
	for _, r := range BuildFakeRecipeList().Recipes {
		recipeIDs = append(recipeIDs, r.ID)
	}

	return &types.MealDatabaseCreationInput{
		ID:            meal.ID,
		Name:          meal.Name,
		Description:   meal.Description,
		CreatedByUser: meal.CreatedByUser,
		Recipes:       recipeIDs,
	}
}
