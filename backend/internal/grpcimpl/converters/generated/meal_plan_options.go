package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealPlanOptionCreationRequestInputToMealPlanOption(input *messages.MealPlanOptionCreationRequestInput) *messages.MealPlanOption {

output := &messages.MealPlanOption{
    AssignedCook: input.AssignedCook,
    MealScale: input.MealScale,
    Notes: input.Notes,
    AssignedDishwasher: input.AssignedDishwasher,
}

return output
}
func ConvertMealPlanOptionUpdateRequestInputToMealPlanOption(input *messages.MealPlanOptionUpdateRequestInput) *messages.MealPlanOption {

output := &messages.MealPlanOption{
    BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
    AssignedDishwasher: input.AssignedDishwasher,
    AssignedCook: input.AssignedCook,
    MealScale: input.MealScale,
    Notes: input.Notes,
}

return output
}
