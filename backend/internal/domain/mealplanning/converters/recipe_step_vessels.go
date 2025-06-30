package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput creates a RecipeStepVesselUpdateRequestInput from a RecipeStepVessel.
func ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(input *mealplanning.RecipeStepVessel) *mealplanning.RecipeStepVesselUpdateRequestInput {
	x := &mealplanning.RecipeStepVesselUpdateRequestInput{
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
func ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(input *mealplanning.RecipeStepVesselCreationRequestInput) *mealplanning.RecipeStepVesselDatabaseCreationInput {
	x := &mealplanning.RecipeStepVesselDatabaseCreationInput{
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
func ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(input *mealplanning.RecipeStepVessel) *mealplanning.RecipeStepVesselCreationRequestInput {
	var instrumentID *string
	if input.Vessel != nil {
		instrumentID = &input.Vessel.ID
	}

	return &mealplanning.RecipeStepVesselCreationRequestInput{
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
func ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(input *mealplanning.RecipeStepVessel) *mealplanning.RecipeStepVesselDatabaseCreationInput {
	var instrumentID *string
	if input.Vessel != nil {
		instrumentID = &input.Vessel.ID
	}

	return &mealplanning.RecipeStepVesselDatabaseCreationInput{
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
