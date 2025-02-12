package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertUserIngredientPreferenceCreationRequestInputToUserIngredientPreference(input *messages.UserIngredientPreferenceCreationRequestInput) *messages.UserIngredientPreference {

output := &messages.UserIngredientPreference{
    Rating: input.Rating,
    Allergy: input.Allergy,
    Notes: input.Notes,
}

return output
}
func ConvertUserIngredientPreferenceUpdateRequestInputToUserIngredientPreference(input *messages.UserIngredientPreferenceUpdateRequestInput) *messages.UserIngredientPreference {

output := &messages.UserIngredientPreference{
    Allergy: input.Allergy,
    Notes: input.Notes,
    Rating: input.Rating,
}

return output
}
