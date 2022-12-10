package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput creates a RecipeStepInstrumentUpdateRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(input *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	x := &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        &input.Instrument.ID,
		Notes:               &input.Notes,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		PreferenceRank:      &input.PreferenceRank,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Optional:            &input.Optional,
		OptionIndex:         &input.OptionIndex,
		MinimumQuantity:     &input.MinimumQuantity,
		MaximumQuantity:     &input.MaximumQuantity,
	}

	return x
}

// ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput creates a RecipeStepInstrumentDatabaseCreationInput from a RecipeStepInstrumentCreationRequestInput.
func ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(input *types.RecipeStepInstrumentCreationRequestInput) *types.RecipeStepInstrumentDatabaseCreationInput {
	x := &types.RecipeStepInstrumentDatabaseCreationInput{
		ID:                  identifiers.New(),
		InstrumentID:        input.InstrumentID,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                input.Name,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
	}

	return x
}

// ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput builds a RecipeStepInstrumentCreationRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentCreationRequestInput {
	var instrumentID *string
	if recipeStepInstrument.Instrument != nil {
		instrumentID = &recipeStepInstrument.Instrument.ID
	}

	return &types.RecipeStepInstrumentCreationRequestInput{
		InstrumentID:        instrumentID,
		Name:                recipeStepInstrument.Name,
		RecipeStepProductID: recipeStepInstrument.RecipeStepProductID,
		Notes:               recipeStepInstrument.Notes,
		PreferenceRank:      recipeStepInstrument.PreferenceRank,
		Optional:            recipeStepInstrument.Optional,
		OptionIndex:         recipeStepInstrument.OptionIndex,
		MinimumQuantity:     recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:     recipeStepInstrument.MaximumQuantity,
	}
}

// ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput builds a RecipeStepInstrumentDatabaseCreationInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(recipeStepInstrument *types.RecipeStepInstrument) *types.RecipeStepInstrumentDatabaseCreationInput {
	var instrumentID *string
	if recipeStepInstrument.Instrument != nil {
		instrumentID = &recipeStepInstrument.Instrument.ID
	}

	return &types.RecipeStepInstrumentDatabaseCreationInput{
		ID:                  recipeStepInstrument.ID,
		InstrumentID:        instrumentID,
		Name:                recipeStepInstrument.Name,
		RecipeStepProductID: recipeStepInstrument.RecipeStepProductID,
		Notes:               recipeStepInstrument.Notes,
		PreferenceRank:      recipeStepInstrument.PreferenceRank,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
		Optional:            recipeStepInstrument.Optional,
		OptionIndex:         recipeStepInstrument.OptionIndex,
		MinimumQuantity:     recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:     recipeStepInstrument.MaximumQuantity,
	}
}
