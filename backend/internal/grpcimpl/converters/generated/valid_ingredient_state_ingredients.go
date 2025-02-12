package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidIngredientStateIngredientCreationRequestInputToValidIngredientStateIngredient(input *messages.ValidIngredientStateIngredientCreationRequestInput) *messages.ValidIngredientStateIngredient {

output := &messages.ValidIngredientStateIngredient{
    Notes: input.Notes,
}

return output
}
func ConvertValidIngredientStateIngredientUpdateRequestInputToValidIngredientStateIngredient(input *messages.ValidIngredientStateIngredientUpdateRequestInput) *messages.ValidIngredientStateIngredient {

output := &messages.ValidIngredientStateIngredient{
    Notes: input.Notes,
}

return output
}
