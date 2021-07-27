package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	return &types.RecipeStep{
		ID:                        uint64(fake.Uint32()),
		ExternalID:                fake.UUID(),
		Index:                     uint(fake.Uint32()),
		PreparationID:             uint64(fake.Uint32()),
		PrerequisiteStep:          uint64(fake.Uint32()),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
		Notes:                     fake.Word(),
		Why:                       fake.Word(),
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe:           fake.Uint64(),
		Ingredients:               BuildFakeRecipeStepIngredientList().RecipeStepIngredients,
	}
}

// BuildFakeRecipeStepList builds a faked RecipeStepList.
func BuildFakeRecipeStepList() *types.RecipeStepList {
	var examples []*types.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStep())
	}

	return &types.RecipeStepList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeSteps: examples,
	}
}

// BuildFakeRecipeStepUpdateInput builds a faked RecipeStepUpdateInput from a recipe step.
func BuildFakeRecipeStepUpdateInput() *types.RecipeStepUpdateInput {
	recipeStep := BuildFakeRecipeStep()
	return &types.RecipeStepUpdateInput{
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Why:                       recipeStep.Why,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepUpdateInputFromRecipeStep builds a faked RecipeStepUpdateInput from a recipe step.
func BuildFakeRecipeStepUpdateInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepUpdateInput {
	return &types.RecipeStepUpdateInput{
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Why:                       recipeStep.Why,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepCreationInput builds a faked RecipeStepCreationInput.
func BuildFakeRecipeStepCreationInput() *types.RecipeStepCreationInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepCreationInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepCreationInputFromRecipeStep builds a faked RecipeStepCreationInput from a recipe step.
func BuildFakeRecipeStepCreationInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepCreationInput {
	ingredients := []*types.RecipeStepIngredientCreationInput{}
	for _, i := range recipeStep.Ingredients {
		ingredients = append(ingredients, BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(i))
	}

	return &types.RecipeStepCreationInput{
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Why:                       recipeStep.Why,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
		Ingredients:               ingredients,
	}
}
