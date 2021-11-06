package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	var ingredients []*types.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		ingredients = append(ingredients, BuildFakeRecipeStepIngredient())
	}

	return &types.RecipeStep{
		ID:                        ksuid.New().String(),
		Index:                     uint(fake.Uint32()),
		Preparation:               *BuildFakeValidPreparation(),
		PrerequisiteStep:          uint64(fake.Uint32()),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
		Notes:                     fake.LoremIpsumSentence(exampleQuantity),
		Why:                       fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe:           ksuid.New().String(),
		Ingredients:               ingredients,
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
			FilteredCount: exampleQuantity,
			TotalCount:    exampleQuantity * 2,
		},
		RecipeSteps: examples,
	}
}

// BuildFakeRecipeStepUpdateRequestInput builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInput() *types.RecipeStepUpdateRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return &types.RecipeStepUpdateRequestInput{
		Index:                     recipeStep.Index,
		Preparation:               recipeStep.Preparation,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Why:                       recipeStep.Why,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepUpdateRequestInputFromRecipeStep builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	return &types.RecipeStepUpdateRequestInput{
		Why:                       recipeStep.Why,
		Index:                     recipeStep.Index,
		Preparation:               recipeStep.Preparation,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepCreationRequestInput builds a faked RecipeStepCreationRequestInput.
func BuildFakeRecipeStepCreationRequestInput() *types.RecipeStepCreationRequestInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepCreationRequestInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepCreationRequestInputFromRecipeStep builds a faked RecipeStepCreationRequestInput from a recipe step.
func BuildFakeRecipeStepCreationRequestInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepCreationRequestInput {
	ingredients := []*types.RecipeStepIngredientCreationRequestInput{}
	for _, ingredient := range recipeStep.Ingredients {
		ingredients = append(ingredients, BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(ingredient))
	}

	return &types.RecipeStepCreationRequestInput{
		ID:                        recipeStep.ID,
		Why:                       recipeStep.Why,
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.Preparation.ID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
		Ingredients:               ingredients,
	}
}

// BuildFakeRecipeStepDatabaseCreationInput builds a faked RecipeStepDatabaseCreationInput.
func BuildFakeRecipeStepDatabaseCreationInput() *types.RecipeStepDatabaseCreationInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep builds a faked RecipeStepDatabaseCreationInput from a recipe step.
func BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepDatabaseCreationInput {
	ingredients := []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, i := range recipeStep.Ingredients {
		ingredients = append(ingredients, BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(i))
	}

	return &types.RecipeStepDatabaseCreationInput{
		ID:                        recipeStep.ID,
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.Preparation.ID,
		Why:                       recipeStep.Why,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Ingredients:               ingredients,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}
