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
		VesselPreposition:    &input.VesselPreposition,
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
		VesselPreposition:               input.VesselPreposition,
		UnavailableAfterStep:            input.UnavailableAfterStep,
	}

	return x
}

// ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput builds a RecipeStepVesselCreationRequestInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(input *types.RecipeStepVessel) *types.RecipeStepVesselCreationRequestInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &types.RecipeStepVesselCreationRequestInput{
		InstrumentID:         instrumentID,
		Name:                 input.Name,
		RecipeStepProductID:  input.RecipeStepProductID,
		Notes:                input.Notes,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		MinimumQuantity:      input.MinimumQuantity,
		MaximumQuantity:      input.MaximumQuantity,
	}
}

// ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput builds a RecipeStepVesselDatabaseCreationInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(input *types.RecipeStepVessel) *types.RecipeStepVesselDatabaseCreationInput {
	var instrumentID *string
	if input.Instrument != nil {
		instrumentID = &input.Instrument.ID
	}

	return &types.RecipeStepVesselDatabaseCreationInput{
		ID:                   input.ID,
		InstrumentID:         instrumentID,
		Name:                 input.Name,
		RecipeStepProductID:  input.RecipeStepProductID,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		MinimumQuantity:      input.MinimumQuantity,
		MaximumQuantity:      input.MaximumQuantity,
	}
}
