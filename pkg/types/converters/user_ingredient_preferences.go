package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput creates a UserIngredientPreferenceUpdateRequestInput from a UserIngredientPreference.
func ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(input *types.UserIngredientPreference) *types.UserIngredientPreferenceUpdateRequestInput {
	x := &types.UserIngredientPreferenceUpdateRequestInput{
		Notes:        &input.Notes,
		IngredientID: &input.Ingredient.ID,
		Rating:       &input.Rating,
		Allergy:      &input.Allergy,
	}

	return x
}

// ConvertUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceDatabaseCreationInput creates a UserIngredientPreferenceDatabaseCreationInput from a UserIngredientPreferenceCreationRequestInput.
func ConvertUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceDatabaseCreationInput(input *types.UserIngredientPreferenceCreationRequestInput) *types.UserIngredientPreferenceDatabaseCreationInput {
	x := &types.UserIngredientPreferenceDatabaseCreationInput{
		ID:           identifiers.New(),
		IngredientID: input.IngredientID,
		Rating:       input.Rating,
		Notes:        input.Notes,
		Allergy:      input.Allergy,
	}

	return x
}

// ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput builds a UserIngredientPreferenceCreationRequestInput from a UserIngredientPreference.
func ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(x *types.UserIngredientPreference) *types.UserIngredientPreferenceCreationRequestInput {
	return &types.UserIngredientPreferenceCreationRequestInput{
		IngredientID: x.Ingredient.ID,
		Rating:       x.Rating,
		Notes:        x.Notes,
		Allergy:      x.Allergy,
	}
}

// ConvertUserIngredientPreferenceToUserIngredientPreferenceDatabaseCreationInput builds a UserIngredientPreferenceDatabaseCreationInput from a UserIngredientPreference.
func ConvertUserIngredientPreferenceToUserIngredientPreferenceDatabaseCreationInput(x *types.UserIngredientPreference) *types.UserIngredientPreferenceDatabaseCreationInput {
	return &types.UserIngredientPreferenceDatabaseCreationInput{
		ID:           x.ID,
		IngredientID: x.Ingredient.ID,
		Rating:       x.Rating,
		Notes:        x.Notes,
		Allergy:      x.Allergy,
	}
}
