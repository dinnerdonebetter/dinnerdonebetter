package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealCreationRequestInputToMeal(input *messages.MealCreationRequestInput) *messages.Meal {
convertedcomponents := make([]*messages.MealComponent, 0, len(input.Components))
for _, item := range input.Components {
    convertedcomponents = append(convertedcomponents, ConvertMealComponentCreationRequestInputToMealComponent(item))
}

output := &messages.Meal{
    EstimatedPortions: input.EstimatedPortions,
    Description: input.Description,
    Name: input.Name,
    Components: convertedcomponents,
    EligibleForMealPlans: input.EligibleForMealPlans,
}

return output
}
func ConvertMealUpdateRequestInputToMeal(input *messages.MealUpdateRequestInput) *messages.Meal {
convertedcomponents := make([]*messages.MealComponent, 0, len(input.Components))
for _, item := range input.Components {
    convertedcomponents = append(convertedcomponents, ConvertMealComponentUpdateRequestInputToMealComponent(item))
}

output := &messages.Meal{
    Components: convertedcomponents,
    EligibleForMealPlans: input.EligibleForMealPlans,
    EstimatedPortions: ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input.EstimatedPortions),
    Description: input.Description,
    CreatedByUser: input.CreatedByUser,
    Name: input.Name,
}

return output
}
