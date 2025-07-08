package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	recipeenumfakes "github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	"github.com/dinnerdonebetter/backend/internal/domain/recipes"
	"github.com/dinnerdonebetter/backend/internal/domain/recipes/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *recipes.RecipeStep {
	recipeStepID := BuildFakeID()

	var ingredients []*recipes.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepIngredient()
		ing.BelongsToRecipeStep = recipeStepID
		ingredients = append(ingredients, ing)
	}

	var instruments []*recipes.RecipeStepInstrument
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepInstrument()
		ing.BelongsToRecipeStep = recipeStepID
		instruments = append(instruments, ing)
	}

	var vessels []*recipes.RecipeStepVessel
	for i := 0; i < exampleQuantity; i++ {
		ing := BuildFakeRecipeStepVessel()
		ing.BelongsToRecipeStep = recipeStepID
		vessels = append(vessels, ing)
	}

	var products []*recipes.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		p := BuildFakeRecipeStepProduct()
		p.BelongsToRecipeStep = recipeStepID
		products = append(products, p)
	}

	completionConditionID := BuildFakeID()
	completionConditions := []*recipes.RecipeStepCompletionCondition{
		{
			ID:                  completionConditionID,
			BelongsToRecipeStep: recipeStepID,
			IngredientState:     recipeenums.ValidIngredientState{},
			Notes:               buildUniqueString(),
			Ingredients: []*recipes.RecipeStepCompletionConditionIngredient{
				{
					ID:                                     BuildFakeID(),
					BelongsToRecipeStepCompletionCondition: completionConditionID,
					RecipeStepIngredient:                   ingredients[0].ID,
				},
			},
			Optional: false,
		},
	}

	return &recipes.RecipeStep{
		ID:                      recipeStepID,
		Index:                   fake.Uint32(),
		Preparation:             *recipeenumfakes.BuildFakeValidPreparation(),
		EstimatedTimeInSeconds:  BuildFakeOptionalUint32Range(),
		TemperatureInCelsius:    BuildFakeOptionalFloat32Range(),
		Notes:                   buildUniqueString(),
		Products:                products,
		Optional:                false,
		CreatedAt:               BuildFakeTime(),
		BelongsToRecipe:         BuildFakeID(),
		Ingredients:             ingredients,
		ExplicitInstructions:    buildUniqueString(),
		ConditionExpression:     buildUniqueString(),
		Instruments:             instruments,
		Vessels:                 vessels,
		CompletionConditions:    completionConditions,
		StartTimerAutomatically: fake.Bool(),
	}
}

// BuildFakeRecipeStepsList builds a faked RecipeStepList.
func BuildFakeRecipeStepsList() *filtering.QueryFilteredResult[recipes.RecipeStep] {
	var examples []*recipes.RecipeStep
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStep())
	}

	return &filtering.QueryFilteredResult[recipes.RecipeStep]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepUpdateRequestInput builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInput() *recipes.RecipeStepUpdateRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepUpdateRequestInput(recipeStep)
}

// BuildFakeRecipeStepCreationRequestInput builds a faked RecipeStepCreationRequestInput.
func BuildFakeRecipeStepCreationRequestInput() *recipes.RecipeStepCreationRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep)
}
