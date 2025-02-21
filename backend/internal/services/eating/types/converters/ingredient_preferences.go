package converters

import (
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

// ConvertIngredientPreferenceToIngredientPreferenceUpdateRequestInput creates a IngredientPreferenceUpdateRequestInput from a IngredientPreferences.
func ConvertIngredientPreferenceToIngredientPreferenceUpdateRequestInput(input *types.IngredientPreference) *types.IngredientPreferenceUpdateRequestInput {
	x := &types.IngredientPreferenceUpdateRequestInput{
		Notes:        &input.Notes,
		IngredientID: &input.Ingredient.ID,
		Rating:       &input.Rating,
		Allergy:      &input.Allergy,
	}

	return x
}

// ConvertIngredientPreferenceCreationRequestInputToIngredientPreferenceDatabaseCreationInput creates a IngredientPreferenceDatabaseCreationInput from a IngredientPreferenceCreationRequestInput.
func ConvertIngredientPreferenceCreationRequestInputToIngredientPreferenceDatabaseCreationInput(input *types.IngredientPreferenceCreationRequestInput) *types.IngredientPreferenceDatabaseCreationInput {
	x := &types.IngredientPreferenceDatabaseCreationInput{
		ValidIngredientGroupID: input.ValidIngredientGroupID,
		ValidIngredientID:      input.ValidIngredientID,
		Rating:                 input.Rating,
		Notes:                  input.Notes,
		Allergy:                input.Allergy,
	}

	return x
}

// ConvertIngredientPreferenceToIngredientPreferenceCreationRequestInput builds a IngredientPreferenceCreationRequestInput from a IngredientPreferences.
func ConvertIngredientPreferenceToIngredientPreferenceCreationRequestInput(x *types.IngredientPreference) *types.IngredientPreferenceCreationRequestInput {
	return &types.IngredientPreferenceCreationRequestInput{
		ValidIngredientID: x.Ingredient.ID,
		Rating:            x.Rating,
		Notes:             x.Notes,
		Allergy:           x.Allergy,
	}
}

// ConvertIngredientPreferenceToIngredientPreferenceDatabaseCreationInput builds a IngredientPreferenceDatabaseCreationInput from a IngredientPreferences.
func ConvertIngredientPreferenceToIngredientPreferenceDatabaseCreationInput(x *types.IngredientPreference) *types.IngredientPreferenceDatabaseCreationInput {
	return &types.IngredientPreferenceDatabaseCreationInput{
		ValidIngredientID: x.Ingredient.ID,
		Rating:            x.Rating,
		Notes:             x.Notes,
		Allergy:           x.Allergy,
		BelongsToUser:     x.BelongsToUser,
	}
}
