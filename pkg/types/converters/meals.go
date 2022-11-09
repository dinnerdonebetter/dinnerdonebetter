package converters

import (
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertMealCreationRequestInputToMealDatabaseCreationInput creates a MealDatabaseCreationInput from a MealCreationRequestInput.
func ConvertMealCreationRequestInputToMealDatabaseCreationInput(input *types.MealCreationRequestInput) *types.MealDatabaseCreationInput {
	x := &types.MealDatabaseCreationInput{
		ID:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		CreatedByUser: input.CreatedByUser,
		Recipes:       input.Recipes,
	}

	return x
}
