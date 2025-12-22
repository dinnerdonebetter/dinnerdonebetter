package recipevalidator

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
)

// RecipeValidator validates and populates recipe creation inputs using bridge table data.
type RecipeValidator struct {
	validIngredientPreparations     map[string]*mealplanning.ValidIngredientPreparation
	validIngredientMeasurementUnits map[string]*mealplanning.ValidIngredientMeasurementUnit
	validPreparationInstruments     map[string]*mealplanning.ValidPreparationInstrument
	validPreparationVessels         map[string]*mealplanning.ValidPreparationVessel
}

// NewRecipeValidator creates a new RecipeValidator with the provided bridge table maps.
func NewRecipeValidator(
	validIngredientPreparations map[string]*mealplanning.ValidIngredientPreparation,
	validIngredientMeasurementUnits map[string]*mealplanning.ValidIngredientMeasurementUnit,
	validPreparationInstruments map[string]*mealplanning.ValidPreparationInstrument,
	validPreparationVessels map[string]*mealplanning.ValidPreparationVessel,
) *RecipeValidator {
	return &RecipeValidator{
		validIngredientPreparations:     validIngredientPreparations,
		validIngredientMeasurementUnits: validIngredientMeasurementUnits,
		validPreparationInstruments:     validPreparationInstruments,
		validPreparationVessels:         validPreparationVessels,
	}
}

// ValidateAndPopulate validates a recipe creation input and populates derived fields.
// It returns an error if any bridge table ID is invalid or doesn't match the step's preparation.
func (v *RecipeValidator) ValidateAndPopulate(input *mealplanning.RecipeDatabaseCreationInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	}

	for stepIdx, step := range input.Steps {
		if err := v.validateStep(stepIdx, step); err != nil {
			return err
		}
	}

	return nil
}

// validateStep validates all ingredients, instruments, and vessels in a step.
func (v *RecipeValidator) validateStep(stepIdx int, step *mealplanning.RecipeStepDatabaseCreationInput) error {
	if step == nil {
		return nil
	}

	// Validate and populate ingredients
	for ingredientIdx, ingredient := range step.Ingredients {
		if err := v.validateAndPopulateIngredient(stepIdx, ingredientIdx, step.PreparationID, ingredient); err != nil {
			return err
		}
	}

	// Validate and populate instruments
	for instrumentIdx, instrument := range step.Instruments {
		if err := v.validateAndPopulateInstrument(stepIdx, instrumentIdx, step.PreparationID, instrument); err != nil {
			return err
		}
	}

	// Validate and populate vessels
	for vesselIdx, vessel := range step.Vessels {
		if err := v.validateAndPopulateVessel(stepIdx, vesselIdx, step.PreparationID, vessel); err != nil {
			return err
		}
	}

	return nil
}

// isRecipeStepProduct returns true if the ingredient is a recipe step product (output from a previous step).
func isRecipeStepProductIngredient(ingredient *mealplanning.RecipeStepIngredientDatabaseCreationInput) bool {
	if ingredient == nil {
		return false
	}
	// A recipe step product is indicated by having a RecipeStepProductID or ProductOfRecipeStepIndex set
	return (ingredient.RecipeStepProductID != nil && *ingredient.RecipeStepProductID != "") ||
		ingredient.ProductOfRecipeStepIndex != nil
}

// isRecipeStepProductInstrument returns true if the instrument is a recipe step product (output from a previous step).
func isRecipeStepProductInstrument(instrument *mealplanning.RecipeStepInstrumentDatabaseCreationInput) bool {
	if instrument == nil {
		return false
	}
	return (instrument.RecipeStepProductID != nil && *instrument.RecipeStepProductID != "") ||
		instrument.ProductOfRecipeStepIndex != nil
}

// isRecipeStepProductVessel returns true if the vessel is a recipe step product (output from a previous step).
func isRecipeStepProductVessel(vessel *mealplanning.RecipeStepVesselDatabaseCreationInput) bool {
	if vessel == nil {
		return false
	}
	return (vessel.RecipeStepProductID != nil && *vessel.RecipeStepProductID != "") ||
		vessel.ProductOfRecipeStepIndex != nil
}

// validateAndPopulateIngredient validates the bridge table IDs for an ingredient and populates derived fields.
func (v *RecipeValidator) validateAndPopulateIngredient(
	stepIdx, ingredientIdx int,
	preparationID string,
	ingredient *mealplanning.RecipeStepIngredientDatabaseCreationInput,
) error {
	if ingredient == nil || isRecipeStepProductIngredient(ingredient) {
		return nil
	}

	// Validate ValidIngredientPreparationID if provided
	if ingredient.ValidIngredientPreparationID != nil && *ingredient.ValidIngredientPreparationID != "" {
		vipID := *ingredient.ValidIngredientPreparationID
		vip, ok := v.validIngredientPreparations[vipID]
		if !ok {
			return fmt.Errorf("step %d ingredient %d: ValidIngredientPreparation %q not found", stepIdx, ingredientIdx, vipID)
		}

		// Validate that the VIP's preparation matches the step's preparation
		if vip.Preparation.ID != preparationID {
			return fmt.Errorf("step %d ingredient %d: ValidIngredientPreparation %q is for preparation %q, but step uses preparation %q",
				stepIdx, ingredientIdx, vipID, vip.Preparation.ID, preparationID)
		}
		ingredient.IngredientID = &vip.Ingredient.ID
	}

	// Validate ValidIngredientMeasurementUnitID if provided
	if ingredient.ValidIngredientMeasurementUnitID != nil && *ingredient.ValidIngredientMeasurementUnitID != "" {
		vimuID := *ingredient.ValidIngredientMeasurementUnitID
		vimu, ok := v.validIngredientMeasurementUnits[vimuID]
		if !ok {
			return fmt.Errorf("step %d ingredient %d: ValidIngredientMeasurementUnit %q not found", stepIdx, ingredientIdx, vimuID)
		}

		// If we have a VIP, validate that the VIMU's ingredient matches
		if ingredient.IngredientID != nil && *ingredient.IngredientID != "" {
			if vimu.Ingredient.ID != *ingredient.IngredientID {
				return fmt.Errorf("step %d ingredient %d: ValidIngredientMeasurementUnit %q is for ingredient %q, but ingredient %q was specified",
					stepIdx, ingredientIdx, vimuID, vimu.Ingredient.ID, *ingredient.IngredientID)
			}
		} else {
			// If no VIP was provided but VIMU was, populate IngredientID from VIMU
			ingredient.IngredientID = &vimu.Ingredient.ID
		}

		ingredient.MeasurementUnitID = vimu.MeasurementUnit.ID
	}

	return nil
}

// validateAndPopulateInstrument validates the bridge table ID for an instrument and populates derived fields.
func (v *RecipeValidator) validateAndPopulateInstrument(
	stepIdx, instrumentIdx int,
	preparationID string,
	instrument *mealplanning.RecipeStepInstrumentDatabaseCreationInput,
) error {
	if instrument == nil || isRecipeStepProductInstrument(instrument) {
		return nil
	}

	// Validate ValidPreparationInstrumentID if provided
	if instrument.ValidPreparationInstrumentID != nil && *instrument.ValidPreparationInstrumentID != "" {
		vpiID := *instrument.ValidPreparationInstrumentID
		vpi, ok := v.validPreparationInstruments[vpiID]
		if !ok {
			return fmt.Errorf("step %d instrument %d: ValidPreparationInstrument %q not found", stepIdx, instrumentIdx, vpiID)
		}

		// Validate that the VPI's preparation matches the step's preparation
		if vpi.Preparation.ID != preparationID {
			return fmt.Errorf("step %d instrument %d: ValidPreparationInstrument %q is for preparation %q, but step uses preparation %q",
				stepIdx, instrumentIdx, vpiID, vpi.Preparation.ID, preparationID)
		}

		// Populate InstrumentID from the VPI
		instrument.InstrumentID = &vpi.Instrument.ID
	}

	return nil
}

// validateAndPopulateVessel validates the bridge table ID for a vessel and populates derived fields.
func (v *RecipeValidator) validateAndPopulateVessel(
	stepIdx, vesselIdx int,
	preparationID string,
	vessel *mealplanning.RecipeStepVesselDatabaseCreationInput,
) error {
	if vessel == nil || isRecipeStepProductVessel(vessel) {
		return nil
	}

	// Validate ValidPreparationVesselID if provided
	if vessel.ValidPreparationVesselID != nil && *vessel.ValidPreparationVesselID != "" {
		vpvID := *vessel.ValidPreparationVesselID
		vpv, ok := v.validPreparationVessels[vpvID]
		if !ok {
			return fmt.Errorf("step %d vessel %d: ValidPreparationVessel %q not found", stepIdx, vesselIdx, vpvID)
		}

		// Validate that the VPV's preparation matches the step's preparation
		if vpv.Preparation.ID != preparationID {
			return fmt.Errorf("step %d vessel %d: ValidPreparationVessel %q is for preparation %q, but step uses preparation %q",
				stepIdx, vesselIdx, vpvID, vpv.Preparation.ID, preparationID)
		}

		// Populate VesselID from the VPV
		vessel.VesselID = &vpv.Vessel.ID
	}

	return nil
}
