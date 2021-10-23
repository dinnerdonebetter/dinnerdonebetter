package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeRecipeStep builds a faked recipe step.
func BuildFakeRecipeStep() *types.RecipeStep {
	return &types.RecipeStep{
		ID:                        ksuid.New().String(),
		Index:                     uint(fake.Uint32()),
		PreparationID:             fake.LoremIpsumSentence(exampleQuantity),
		PrerequisiteStep:          uint64(fake.Uint32()),
		MinEstimatedTimeInSeconds: fake.Uint32(),
		MaxEstimatedTimeInSeconds: fake.Uint32(),
		TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
		Why:                       fake.LoremIpsumSentence(exampleQuantity),
		Notes:                     fake.LoremIpsumSentence(exampleQuantity),
		RecipeID:                  fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
		BelongsToRecipe:           fake.UUID(),
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
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		Why:                       recipeStep.Why,
		RecipeID:                  recipeStep.RecipeID,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepUpdateRequestInputFromRecipeStep builds a faked RecipeStepUpdateRequestInput from a recipe step.
func BuildFakeRecipeStepUpdateRequestInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	return &types.RecipeStepUpdateRequestInput{
		Why:                       recipeStep.Why,
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		RecipeID:                  recipeStep.RecipeID,
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
	return &types.RecipeStepCreationRequestInput{
		ID:                        recipeStep.ID,
		Why:                       recipeStep.Why,
		Index:                     recipeStep.Index,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		RecipeID:                  recipeStep.RecipeID,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}

// BuildFakeRecipeStepDatabaseCreationInput builds a faked RecipeStepDatabaseCreationInput.
func BuildFakeRecipeStepDatabaseCreationInput() *types.RecipeStepDatabaseCreationInput {
	recipeStep := BuildFakeRecipeStep()
	return BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(recipeStep)
}

// BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep builds a faked RecipeStepDatabaseCreationInput from a recipe step.
func BuildFakeRecipeStepDatabaseCreationInputFromRecipeStep(recipeStep *types.RecipeStep) *types.RecipeStepDatabaseCreationInput {
	return &types.RecipeStepDatabaseCreationInput{
		ID:                        recipeStep.ID,
		Index:                     recipeStep.Index,
		Why:                       recipeStep.Why,
		PreparationID:             recipeStep.PreparationID,
		PrerequisiteStep:          recipeStep.PrerequisiteStep,
		MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
		Notes:                     recipeStep.Notes,
		RecipeID:                  recipeStep.RecipeID,
		BelongsToRecipe:           recipeStep.BelongsToRecipe,
	}
}
