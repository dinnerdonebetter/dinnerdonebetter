package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/recipes"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput creates a RecipeStepVesselUpdateRequestInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(input *recipes.RecipeStepVessel) *recipes.RecipeStepVesselUpdateRequestInput {
	x := &recipes.RecipeStepVesselUpdateRequestInput{
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
func ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(input *recipes.RecipeStepVesselCreationRequestInput) *recipes.RecipeStepVesselDatabaseCreationInput {
	x := &recipes.RecipeStepVesselDatabaseCreationInput{
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
func ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(input *recipes.RecipeStepVessel) *recipes.RecipeStepVesselCreationRequestInput {
	var instrumentID *string
	if input.Vessel != nil {
		instrumentID = &input.Vessel.ID
	}

	return &recipes.RecipeStepVesselCreationRequestInput{
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
func ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(input *recipes.RecipeStepVessel) *recipes.RecipeStepVesselDatabaseCreationInput {
	var instrumentID *string
	if input.Vessel != nil {
		instrumentID = &input.Vessel.ID
	}

	return &recipes.RecipeStepVesselDatabaseCreationInput{
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
