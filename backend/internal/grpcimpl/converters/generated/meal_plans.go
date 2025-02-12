package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealPlanCreationRequestInputToMealPlan(input *messages.MealPlanCreationRequestInput) *messages.MealPlan {
convertedevents := make([]*messages.MealPlanEvent, 0, len(input.Events))
for _, item := range input.Events {
    convertedevents = append(convertedevents, ConvertMealPlanEventCreationRequestInputToMealPlanEvent(item))
}

output := &messages.MealPlan{
    VotingDeadline: input.VotingDeadline,
    ElectionMethod: input.ElectionMethod,
    Notes: input.Notes,
    Events: convertedevents,
}

return output
}
func ConvertMealPlanUpdateRequestInputToMealPlan(input *messages.MealPlanUpdateRequestInput) *messages.MealPlan {

output := &messages.MealPlan{
    VotingDeadline: input.VotingDeadline,
    Notes: input.Notes,
    BelongsToHousehold: input.BelongsToHousehold,
}

return output
}
