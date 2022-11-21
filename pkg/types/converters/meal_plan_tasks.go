package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertMealPlanTaskCreationRequestInputToMealPlanTaskDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertMealPlanTaskCreationRequestInputToMealPlanTaskDatabaseCreationInput(input *types.MealPlanTaskCreationRequestInput) *types.MealPlanTaskDatabaseCreationInput {
	x := &types.MealPlanTaskDatabaseCreationInput{
		ID:                  identifiers.New(),
		AssignedToUser:      input.AssignedToUser,
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOptionID:    input.MealPlanOptionID,
		RecipePrepTaskID:    input.RecipePrepTaskID,
	}

	return x
}

// ConvertMealPlanTaskToMealPlanTaskCreationRequestInput builds a meal plan task.
func ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(x *types.MealPlanTask) *types.MealPlanTaskCreationRequestInput {
	return &types.MealPlanTaskCreationRequestInput{
		Status:              x.Status,
		StatusExplanation:   x.StatusExplanation,
		CreationExplanation: x.CreationExplanation,
	}
}
