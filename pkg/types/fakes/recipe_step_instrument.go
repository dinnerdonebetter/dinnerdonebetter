package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipeStepInstrument builds a faked recipe step instrument.
func BuildFakeRecipeStepInstrument() *types.RecipeStepInstrument {
	return &types.RecipeStepInstrument{
		ID:                  ksuid.New().String(),
		InstrumentID:        func(x string) *string { return &x }(fake.LoremIpsumSentence(exampleQuantity)),
		RecipeStepID:        fake.LoremIpsumSentence(exampleQuantity),
		Notes:               fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.UUID(),
	}
}

// BuildFakeRecipeStepInstrumentList builds a faked RecipeStepInstrumentList.
func BuildFakeRecipeStepInstrumentList() *types.RecipeStepInstrumentList {
	var examples []*types.RecipeStepInstrument
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepInstrument())
	}

	return &types.RecipeStepInstrumentList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeStepInstruments: examples,
	}
}

// BuildFakeRecipeStepInstrumentUpdateRequestInput builds a faked RecipeStepInstrumentUpdateRequestInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentUpdateRequestInput() *types.RecipeStepInstrumentUpdateRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        recipeStepInstrument.InstrumentID,
		RecipeStepID:        recipeStepInstrument.RecipeStepID,
		Notes:               recipeStepInstrument.Notes,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepInstrumentUpdateRequestInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentUpdateRequestInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentUpdateRequestInputFromRecipeStepInstrument(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	return &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        recipeStepInstrument.InstrumentID,
		RecipeStepID:        recipeStepInstrument.RecipeStepID,
		Notes:               recipeStepInstrument.Notes,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepInstrumentCreationRequestInput builds a faked RecipeStepInstrumentCreationRequestInput.
func BuildFakeRecipeStepInstrumentCreationRequestInput() *types.RecipeStepInstrumentCreationRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(recipeStepInstrument)
}

// BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentCreationRequestInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentCreationRequestInput {
	return &types.RecipeStepInstrumentCreationRequestInput{
		ID:                  recipeStepInstrument.ID,
		InstrumentID:        recipeStepInstrument.InstrumentID,
		RecipeStepID:        recipeStepInstrument.RecipeStepID,
		Notes:               recipeStepInstrument.Notes,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepInstrumentDatabaseCreationInput builds a faked RecipeStepInstrumentDatabaseCreationInput.
func BuildFakeRecipeStepInstrumentDatabaseCreationInput() *types.RecipeStepInstrumentDatabaseCreationInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument(recipeStepInstrument)
}

// BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentDatabaseCreationInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentDatabaseCreationInput {
	return &types.RecipeStepInstrumentDatabaseCreationInput{
		ID:                  recipeStepInstrument.ID,
		InstrumentID:        recipeStepInstrument.InstrumentID,
		RecipeStepID:        recipeStepInstrument.RecipeStepID,
		Notes:               recipeStepInstrument.Notes,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
	}
}
