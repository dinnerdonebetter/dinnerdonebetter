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
		Instrument:          BuildFakeValidInstrument(),
		Name:                buildUniqueString(),
		ProductOfRecipeStep: fake.Bool(),
		RecipeStepProductID: nil,
		Notes:               buildUniqueString(),
		PreferenceRank:      fake.Uint8(),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.UUID(),
		Optional:            fake.Bool(),
		MinimumQuantity:     fake.Uint32(),
		MaximumQuantity:     fake.Uint32(),
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
		InstrumentID:        &recipeStepInstrument.Instrument.ID,
		RecipeStepProductID: recipeStepInstrument.RecipeStepProductID,
		Name:                &recipeStepInstrument.Name,
		ProductOfRecipeStep: &recipeStepInstrument.ProductOfRecipeStep,
		Notes:               &recipeStepInstrument.Notes,
		PreferenceRank:      &recipeStepInstrument.PreferenceRank,
		BelongsToRecipeStep: &recipeStepInstrument.BelongsToRecipeStep,
		Optional:            &recipeStepInstrument.Optional,
		MinimumQuantity:     &recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:     &recipeStepInstrument.MaximumQuantity,
	}
}

// BuildFakeRecipeStepInstrumentUpdateRequestInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentUpdateRequestInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentUpdateRequestInputFromRecipeStepInstrument(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	return &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        &recipeStepInstrument.Instrument.ID,
		Name:                &recipeStepInstrument.Name,
		ProductOfRecipeStep: &recipeStepInstrument.ProductOfRecipeStep,
		RecipeStepProductID: recipeStepInstrument.RecipeStepProductID,
		Notes:               &recipeStepInstrument.Notes,
		PreferenceRank:      &recipeStepInstrument.PreferenceRank,
		BelongsToRecipeStep: &recipeStepInstrument.BelongsToRecipeStep,
		Optional:            &recipeStepInstrument.Optional,
		MinimumQuantity:     &recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:     &recipeStepInstrument.MaximumQuantity,
	}
}

// BuildFakeRecipeStepInstrumentCreationRequestInput builds a faked RecipeStepInstrumentCreationRequestInput.
func BuildFakeRecipeStepInstrumentCreationRequestInput() *types.RecipeStepInstrumentCreationRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(recipeStepInstrument)
}

// BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentCreationRequestInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentCreationRequestInput {
	var instrumentID *string
	if recipeStepInstrument.Instrument != nil {
		instrumentID = &recipeStepInstrument.Instrument.ID
	}

	return &types.RecipeStepInstrumentCreationRequestInput{
		ID:                  recipeStepInstrument.ID,
		InstrumentID:        instrumentID,
		Name:                recipeStepInstrument.Name,
		ProductOfRecipeStep: recipeStepInstrument.ProductOfRecipeStep,
		RecipeStepProductID: recipeStepInstrument.RecipeStepProductID,
		Notes:               recipeStepInstrument.Notes,
		PreferenceRank:      recipeStepInstrument.PreferenceRank,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
		Optional:            recipeStepInstrument.Optional,
		MinimumQuantity:     recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:     recipeStepInstrument.MaximumQuantity,
	}
}

// BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentDatabaseCreationInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentDatabaseCreationInput {
	var instrumentID *string
	if recipeStepInstrument.Instrument != nil {
		instrumentID = &recipeStepInstrument.Instrument.ID
	}

	return &types.RecipeStepInstrumentDatabaseCreationInput{
		ID:                  recipeStepInstrument.ID,
		InstrumentID:        instrumentID,
		Name:                recipeStepInstrument.Name,
		ProductOfRecipeStep: recipeStepInstrument.ProductOfRecipeStep,
		RecipeStepProductID: recipeStepInstrument.RecipeStepProductID,
		Notes:               recipeStepInstrument.Notes,
		PreferenceRank:      recipeStepInstrument.PreferenceRank,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
		Optional:            recipeStepInstrument.Optional,
		MinimumQuantity:     recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:     recipeStepInstrument.MaximumQuantity,
	}
}
