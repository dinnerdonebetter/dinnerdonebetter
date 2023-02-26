package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput creates a RecipeStepVesselUpdateRequestInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(input *types.RecipeStepVessel) *types.RecipeStepVesselUpdateRequestInput {
	x := &types.RecipeStepVesselUpdateRequestInput{
		InstrumentID:         &input.Instrument.ID,
		Notes:                &input.Notes,
		RecipeStepProductID:  input.RecipeStepProductID,
		Name:                 &input.Name,
		BelongsToRecipeStep:  &input.BelongsToRecipeStep,
		MinimumQuantity:      &input.MinimumQuantity,
		MaximumQuantity:      input.MaximumQuantity,
		VesselPredicate:      &input.VesselPredicate,
		UnavailableAfterStep: &input.UnavailableAfterStep,
	}

	return x
}

// ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput creates a RecipeStepVesselDatabaseCreationInput from a RecipeStepVesselCreationRequestInput.
func ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(input *types.RecipeStepVesselCreationRequestInput) *types.RecipeStepVesselDatabaseCreationInput {
	x := &types.RecipeStepVesselDatabaseCreationInput{
		ID:                              identifiers.New(),
		InstrumentID:                    input.InstrumentID,
		RecipeStepProductID:             input.RecipeStepProductID,
		Name:                            input.Name,
		Notes:                           input.Notes,
		MinimumQuantity:                 input.MinimumQuantity,
		MaximumQuantity:                 input.MaximumQuantity,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
		VesselPredicate:                 input.VesselPredicate,
		UnavailableAfterStep:            input.UnavailableAfterStep,
	}

	return x
}

// ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput builds a RecipeStepVesselCreationRequestInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(recipeStepInstrument *types.RecipeStepVessel) *types.RecipeStepVesselCreationRequestInput {
	var instrumentID *string
	if recipeStepInstrument.Instrument != nil {
		instrumentID = &recipeStepInstrument.Instrument.ID
	}

	return &types.RecipeStepVesselCreationRequestInput{
		InstrumentID:         instrumentID,
		Name:                 recipeStepInstrument.Name,
		RecipeStepProductID:  recipeStepInstrument.RecipeStepProductID,
		Notes:                recipeStepInstrument.Notes,
		VesselPredicate:      recipeStepInstrument.VesselPredicate,
		UnavailableAfterStep: recipeStepInstrument.UnavailableAfterStep,
		MinimumQuantity:      recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:      recipeStepInstrument.MaximumQuantity,
	}
}

// ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput builds a RecipeStepVesselDatabaseCreationInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(recipeStepInstrument *types.RecipeStepVessel) *types.RecipeStepVesselDatabaseCreationInput {
	var instrumentID *string
	if recipeStepInstrument.Instrument != nil {
		instrumentID = &recipeStepInstrument.Instrument.ID
	}

	return &types.RecipeStepVesselDatabaseCreationInput{
		ID:                   recipeStepInstrument.ID,
		InstrumentID:         instrumentID,
		Name:                 recipeStepInstrument.Name,
		RecipeStepProductID:  recipeStepInstrument.RecipeStepProductID,
		Notes:                recipeStepInstrument.Notes,
		BelongsToRecipeStep:  recipeStepInstrument.BelongsToRecipeStep,
		VesselPredicate:      recipeStepInstrument.VesselPredicate,
		UnavailableAfterStep: recipeStepInstrument.UnavailableAfterStep,
		MinimumQuantity:      recipeStepInstrument.MinimumQuantity,
		MaximumQuantity:      recipeStepInstrument.MaximumQuantity,
	}
}
