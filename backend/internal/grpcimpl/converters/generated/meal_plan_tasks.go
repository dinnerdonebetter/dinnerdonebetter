package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealPlanTaskCreationRequestInputToMealPlanTask(input *messages.MealPlanTaskCreationRequestInput) *messages.MealPlanTask {

output := &messages.MealPlanTask{
    CreationExplanation: input.CreationExplanation,
    StatusExplanation: input.StatusExplanation,
    AssignedToUser: input.AssignedToUser,
    Status: input.Status,
}

return output
}
