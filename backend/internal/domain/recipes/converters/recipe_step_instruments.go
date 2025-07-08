package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/recipes"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput creates a RecipeStepInstrumentUpdateRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(input *recipes.RecipeStepInstrument) *recipes.RecipeStepInstrumentUpdateRequestInput {
	x := &recipes.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID:        &input.Instrument.ID,
		Notes:               &input.Notes,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		PreferenceRank:      &input.PreferenceRank,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Optional:            &input.Optional,
		OptionIndex:         &input.OptionIndex,
		Quantity: types.Uint32RangeWithOptionalMaxUpdateRequestInput{
			Min: &input.Quantity.Min,
			Max: input.Quantity.Max,
		},
	}

	return x
}

// ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput creates a RecipeStepInstrumentDatabaseCreationInput from a RecipeStepInstrumentCreationRequestInput.
func ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(input *recipes.RecipeStepInstrumentCreationRequestInput) *recipes.RecipeStepInstrumentDatabaseCreationInput {
	x := &recipes.RecipeStepInstrumentDatabaseCreationInput{
		ID:                              identifiers.New(),
		InstrumentID:                    input.InstrumentID,
		RecipeStepProductID:             input.RecipeStepProductID,
		Name:                            input.Name,
		Notes:                           input.Notes,
		PreferenceRank:                  input.PreferenceRank,
		Optional:                        input.Optional,
		OptionIndex:                     input.OptionIndex,
		Quantity:                        input.Quantity,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
	}

	return x
}

// ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput builds a RecipeStepInstrumentCreationRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(input *recipes.RecipeStepInstrument) *recipes.RecipeStepInstrumentCreationRequestInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &recipes.RecipeStepInstrumentCreationRequestInput{
		InstrumentID:        instrumentID,
		Name:                input.Name,
		RecipeStepProductID: input.RecipeStepProductID,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		Quantity:            input.Quantity,
	}
}

// ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput builds a RecipeStepInstrumentDatabaseCreationInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(input *recipes.RecipeStepInstrument) *recipes.RecipeStepInstrumentDatabaseCreationInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &recipes.RecipeStepInstrumentDatabaseCreationInput{
		ID:                  input.ID,
		InstrumentID:        instrumentID,
		Name:                input.Name,
		RecipeStepProductID: input.RecipeStepProductID,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		Quantity:            input.Quantity,
	}
}
