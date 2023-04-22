package fakes

import (
	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepInstrument builds a faked recipe step instrument.
func BuildFakeRecipeStepInstrument() *types.RecipeStepInstrument {
	return &types.RecipeStepInstrument{
		ID:                  BuildFakeID(),
		Instrument:          BuildFakeValidInstrument(),
		Name:                buildUniqueString(),
		RecipeStepProductID: nil,
		Notes:               buildUniqueString(),
		PreferenceRank:      fake.Uint8(),
		CreatedAt:           BuildFakeTime(),
		BelongsToRecipeStep: fake.UUID(),
		Optional:            fake.Bool(),
		OptionIndex:         uint16(fake.Uint8()),
		MinimumQuantity:     fake.Uint32(),
		MaximumQuantity:     pointers.Pointer(fake.Uint32()),
	}
}

// BuildFakeRecipeStepInstrumentList builds a faked RecipeStepInstrumentList.
func BuildFakeRecipeStepInstrumentList() *types.QueryFilteredResult[types.RecipeStepInstrument] {
	var examples []*types.RecipeStepInstrument
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepInstrument())
	}

	return &types.QueryFilteredResult[types.RecipeStepInstrument]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
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
