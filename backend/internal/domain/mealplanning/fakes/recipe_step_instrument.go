package fakes

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	recipeenumfakes "github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStepInstrument builds a faked recipe step instrument.
func BuildFakeRecipeStepInstrument() *types.RecipeStepInstrument {
	return &types.RecipeStepInstrument{
		ID:                  BuildFakeID(),
		Instrument:          recipeenumfakes.BuildFakeValidInstrument(),
		Name:                buildUniqueString(),
		RecipeStepProductID: nil,
		Notes:               buildUniqueString(),
		PreferenceRank:      fake.Uint8(),
		CreatedAt:           BuildFakeTime(),
		BelongsToRecipeStep: fake.UUID(),
		Optional:            fake.Bool(),
		OptionIndex:         uint16(fake.Uint8()),
		Quantity:            BuildFakeUint32RangeWithOptionalMax(),
	}
}

// BuildFakeRecipeStepInstrumentsList builds a faked RecipeStepInstrumentList.
func BuildFakeRecipeStepInstrumentsList() *filtering.QueryFilteredResult[types.RecipeStepInstrument] {
	var examples []*types.RecipeStepInstrument
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepInstrument())
	}

	return &filtering.QueryFilteredResult[types.RecipeStepInstrument]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepInstrumentUpdateRequestInput builds a faked RecipeStepInstrumentUpdateRequestInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentUpdateRequestInput() *types.RecipeStepInstrumentUpdateRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(recipeStepInstrument)
}

// BuildFakeRecipeStepInstrumentCreationRequestInput builds a faked RecipeStepInstrumentCreationRequestInput.
func BuildFakeRecipeStepInstrumentCreationRequestInput() *types.RecipeStepInstrumentCreationRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(recipeStepInstrument)
}
