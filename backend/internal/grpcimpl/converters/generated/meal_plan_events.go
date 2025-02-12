package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealPlanEventCreationRequestInputToMealPlanEvent(input *messages.MealPlanEventCreationRequestInput) *messages.MealPlanEvent {
convertedoptions := make([]*messages.MealPlanOption, 0, len(input.Options))
for _, item := range input.Options {
    convertedoptions = append(convertedoptions, ConvertMealPlanOptionCreationRequestInputToMealPlanOption(item))
}

output := &messages.MealPlanEvent{
    EndsAt: input.EndsAt,
    MealName: input.MealName,
    Notes: input.Notes,
    Options: convertedoptions,
    StartsAt: input.StartsAt,
}

return output
}
func ConvertMealPlanEventUpdateRequestInputToMealPlanEvent(input *messages.MealPlanEventUpdateRequestInput) *messages.MealPlanEvent {

output := &messages.MealPlanEvent{
    EndsAt: input.EndsAt,
    MealName: input.MealName,
    Notes: input.Notes,
    BelongsToMealPlan: input.BelongsToMealPlan,
    StartsAt: input.StartsAt,
}

return output
}
