package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	recipeStepID := BuildFakeID()

	var ingredients []*types.RecipeStepIngredient
	for range exampleQuantity {
		ing := BuildFakeRecipeStepIngredient()
		ing.BelongsToRecipeStep = recipeStepID
		ingredients = append(ingredients, ing)
	}

	var instruments []*types.RecipeStepInstrument
	for range exampleQuantity {
		ing := BuildFakeRecipeStepInstrument()
		ing.BelongsToRecipeStep = recipeStepID
		instruments = append(instruments, ing)
	}

	var vessels []*types.RecipeStepVessel
	for range exampleQuantity {
		ing := BuildFakeRecipeStepVessel()
		ing.BelongsToRecipeStep = recipeStepID
		vessels = append(vessels, ing)
	}

	var products []*types.RecipeStepProduct
	for range exampleQuantity {
		p := BuildFakeRecipeStepProduct()
		p.BelongsToRecipeStep = recipeStepID
		products = append(products, p)
	}

	completionConditionID := BuildFakeID()
	completionConditions := []*types.RecipeStepCompletionCondition{
		{
			ID:                  completionConditionID,
			BelongsToRecipeStep: recipeStepID,
			IngredientState:     types.ValidIngredientState{},
			Notes:               buildUniqueString(),
			Ingredients: []*types.RecipeStepCompletionConditionIngredient{
				{
					ID:                                     BuildFakeID(),
					BelongsToRecipeStepCompletionCondition: completionConditionID,
					RecipeStepIngredient:                   ingredients[0].ID,
				},
			},
			Optional: false,
		},
	}

	return &types.RecipeStep{
		ID:                      recipeStepID,
		Index:                   fake.Uint32(),
		Preparation:             *BuildFakeValidPreparation(),
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
func BuildFakeRecipeStepsList() *filtering.QueryFilteredResult[types.RecipeStep] {
	var examples []*types.RecipeStep
	for range exampleQuantity {
		examples = append(examples, BuildFakeRecipeStep())
	}

	return &filtering.QueryFilteredResult[types.RecipeStep]{
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
func BuildFakeRecipeStepUpdateRequestInput() *types.RecipeStepUpdateRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepUpdateRequestInput(recipeStep)
}

// BuildFakeRecipeStepCreationRequestInput builds a faked RecipeStepCreationRequestInput.
func BuildFakeRecipeStepCreationRequestInput() *types.RecipeStepCreationRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return converters.ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep)
}
