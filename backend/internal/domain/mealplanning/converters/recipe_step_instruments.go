package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput creates a RecipeStepInstrumentUpdateRequestInput from a RecipeStepInstrument.
func ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(input *mealplanning.RecipeStepInstrument) *mealplanning.RecipeStepInstrumentUpdateRequestInput {
	x := &mealplanning.RecipeStepInstrumentUpdateRequestInput{
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
func ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(input *mealplanning.RecipeStepInstrumentCreationRequestInput) *mealplanning.RecipeStepInstrumentDatabaseCreationInput {
	x := &mealplanning.RecipeStepInstrumentDatabaseCreationInput{
		ID:                              identifiers.New(),
		InstrumentID:                    input.InstrumentID,
		ValidPreparationInstrumentID:    input.ValidPreparationInstrumentID,
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
func ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(input *mealplanning.RecipeStepInstrument) *mealplanning.RecipeStepInstrumentCreationRequestInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &mealplanning.RecipeStepInstrumentCreationRequestInput{
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
func ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(input *mealplanning.RecipeStepInstrument) *mealplanning.RecipeStepInstrumentDatabaseCreationInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &mealplanning.RecipeStepInstrumentDatabaseCreationInput{
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
