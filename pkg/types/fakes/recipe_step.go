package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/internal/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	recipeStepID := BuildFakeID()

	minTemp := BuildFakeNumber()

	var ingredients []*types.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepIngredient()
		ing.BelongsToRecipeStep = recipeStepID
		ingredients = append(ingredients, ing)
	}

	var instruments []*types.RecipeStepInstrument
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepInstrument()
		ing.BelongsToRecipeStep = recipeStepID
		instruments = append(instruments, ing)
	}

	var products []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		p := BuildFakeRecipeStepProduct()
		p.BelongsToRecipeStep = recipeStepID
		products = append(products, p)
	}

	return &types.RecipeStep{
		ID:                            recipeStepID,
		Index:                         fake.Uint32(),
		Preparation:                   *BuildFakeValidPreparation(),
		MinimumEstimatedTimeInSeconds: func(x uint32) *uint32 { return &x }(fake.Uint32()),
		MaximumEstimatedTimeInSeconds: func(x uint32) *uint32 { return &x }(fake.Uint32()),
		MinimumTemperatureInCelsius:   pointers.Float32(float32(minTemp)),
		MaximumTemperatureInCelsius:   pointers.Float32(float32(minTemp + 1)),
		Notes:                         buildUniqueString(),
		Products:                      products,
		Optional:                      false,
		CreatedAt:                     fake.Date(),
		BelongsToRecipe:               BuildFakeID(),
		Ingredients:                   ingredients,
		ExplicitInstructions:          buildUniqueString(),
		Instruments:                   instruments,
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

// BuildFakeRecipeStepUpdateRequestInput builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInput() *types.RecipeStepUpdateRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return &types.RecipeStepUpdateRequestInput{
		Index:                         &recipeStep.Index,
		Preparation:                   &recipeStep.Preparation,
		MinimumEstimatedTimeInSeconds: recipeStep.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: recipeStep.MaximumEstimatedTimeInSeconds,
		MaximumTemperatureInCelsius:   recipeStep.MaximumTemperatureInCelsius,
		MinimumTemperatureInCelsius:   recipeStep.MinimumTemperatureInCelsius,
		Notes:                         &recipeStep.Notes,
		Optional:                      &recipeStep.Optional,
		ExplicitInstructions:          &recipeStep.ExplicitInstructions,
		BelongsToRecipe:               recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepCreationRequestInput builds a faked RecipeStepCreationRequestInput.
func BuildFakeRecipeStepCreationRequestInput() *types.RecipeStepCreationRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep)
}
