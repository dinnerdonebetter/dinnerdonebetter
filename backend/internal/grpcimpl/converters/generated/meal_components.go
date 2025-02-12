package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealComponentCreationRequestInputToMealComponent(input *messages.MealComponentCreationRequestInput) *messages.MealComponent {

output := &messages.MealComponent{
    ComponentType: input.ComponentType,
    RecipeScale: input.RecipeScale,
}

return output
}
func ConvertMealComponentUpdateRequestInputToMealComponent(input *messages.MealComponentUpdateRequestInput) *messages.MealComponent {

output := &messages.MealComponent{
    ComponentType: input.ComponentType,
    RecipeScale: input.RecipeScale,
}

return output
}
