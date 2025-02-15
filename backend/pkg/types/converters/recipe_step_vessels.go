package converters

import (
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput creates a RecipeStepVesselUpdateRequestInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(input *types.RecipeStepVessel) *types.RecipeStepVesselUpdateRequestInput {
	x := &types.RecipeStepVesselUpdateRequestInput{
		VesselID:            &input.Vessel.ID,
		Notes:               &input.Notes,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Quantity: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: &input.Quantity.Min,
			Max: input.Quantity.Max,
		},
		VesselPreposition:    &input.VesselPreposition,
		UnavailableAfterStep: &input.UnavailableAfterStep,
	}

	return x
}

// ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput creates a RecipeStepVesselDatabaseCreationInput from a RecipeStepVesselCreationRequestInput.
func ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(input *types.RecipeStepVesselCreationRequestInput) *types.RecipeStepVesselDatabaseCreationInput {
	x := &types.RecipeStepVesselDatabaseCreationInput{
		ID:                              identifiers.New(),
		VesselID:                        input.VesselID,
		RecipeStepProductID:             input.RecipeStepProductID,
		Name:                            input.Name,
		Notes:                           input.Notes,
		Quantity:                        input.Quantity,
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
	if input.Vessel != nil {
		instrumentID = &input.Vessel.ID
	}

	return &types.RecipeStepVesselCreationRequestInput{
		VesselID:             instrumentID,
		Name:                 input.Name,
		RecipeStepProductID:  input.RecipeStepProductID,
		Notes:                input.Notes,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		Quantity:             input.Quantity,
	}
}

// ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput builds a RecipeStepVesselDatabaseCreationInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(input *types.RecipeStepVessel) *types.RecipeStepVesselDatabaseCreationInput {
	var instrumentID *string
	if input.Vessel != nil {
		instrumentID = &input.Vessel.ID
	}

	return &types.RecipeStepVesselDatabaseCreationInput{
		ID:                   input.ID,
		VesselID:             instrumentID,
		Name:                 input.Name,
		RecipeStepProductID:  input.RecipeStepProductID,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		Quantity:             input.Quantity,
	}
}
