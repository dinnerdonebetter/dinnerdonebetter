package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidIngredientPreparationCreationRequestInputToValidIngredientPreparation(input *messages.ValidIngredientPreparationCreationRequestInput) *messages.ValidIngredientPreparation {

output := &messages.ValidIngredientPreparation{
    Notes: input.Notes,
}

return output
}
func ConvertValidIngredientPreparationUpdateRequestInputToValidIngredientPreparation(input *messages.ValidIngredientPreparationUpdateRequestInput) *messages.ValidIngredientPreparation {

output := &messages.ValidIngredientPreparation{
    Notes: input.Notes,
}

return output
}
