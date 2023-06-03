package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertMealRatingToMealRatingUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertMealRatingToMealRatingUpdateRequestInput(x *types.MealRating) *types.MealRatingUpdateRequestInput {
	out := &types.MealRatingUpdateRequestInput{
		MealID:       &x.MealID,
		Taste:        &x.Taste,
		Difficulty:   &x.Difficulty,
		Cleanup:      &x.Cleanup,
		Instructions: &x.Instructions,
		Overall:      &x.Overall,
		Notes:        &x.Notes,
		ByUser:       &x.ByUser,
	}

	return out
}

// ConvertMealRatingCreationRequestInputToMealRatingDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertMealRatingCreationRequestInputToMealRatingDatabaseCreationInput(x *types.MealRatingCreationRequestInput) *types.MealRatingDatabaseCreationInput {
	out := &types.MealRatingDatabaseCreationInput{
		ID:           identifiers.New(),
		MealID:       x.MealID,
		Notes:        x.Notes,
		ByUser:       x.ByUser,
		Taste:        x.Taste,
		Difficulty:   x.Difficulty,
		Cleanup:      x.Cleanup,
		Instructions: x.Instructions,
		Overall:      x.Overall,
	}

	return out
}

// ConvertMealRatingToMealRatingCreationRequestInput builds a MealRatingCreationRequestInput from a Ingredient.
func ConvertMealRatingToMealRatingCreationRequestInput(x *types.MealRating) *types.MealRatingCreationRequestInput {
	return &types.MealRatingCreationRequestInput{
		MealID:       x.MealID,
		Notes:        x.Notes,
		ByUser:       x.ByUser,
		Taste:        x.Taste,
		Difficulty:   x.Difficulty,
		Cleanup:      x.Cleanup,
		Instructions: x.Instructions,
		Overall:      x.Overall,
	}
}

// ConvertMealRatingToMealRatingDatabaseCreationInput builds a MealRatingDatabaseCreationInput from a MealRating.
func ConvertMealRatingToMealRatingDatabaseCreationInput(x *types.MealRating) *types.MealRatingDatabaseCreationInput {
	return &types.MealRatingDatabaseCreationInput{
		ID:           x.ID,
		MealID:       x.MealID,
		Notes:        x.Notes,
		ByUser:       x.ByUser,
		Taste:        x.Taste,
		Difficulty:   x.Difficulty,
		Cleanup:      x.Cleanup,
		Instructions: x.Instructions,
		Overall:      x.Overall,
	}
}
