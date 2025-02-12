package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidIngredientStateCreationRequestInputToValidIngredientState(input *messages.ValidIngredientStateCreationRequestInput) *messages.ValidIngredientState {

output := &messages.ValidIngredientState{
    PastTense: input.PastTense,
    Description: input.Description,
    IconPath: input.IconPath,
    Name: input.Name,
    AttributeType: input.AttributeType,
    Slug: input.Slug,
}

return output
}
func ConvertValidIngredientStateUpdateRequestInputToValidIngredientState(input *messages.ValidIngredientStateUpdateRequestInput) *messages.ValidIngredientState {

output := &messages.ValidIngredientState{
    PastTense: input.PastTense,
    Description: input.Description,
    IconPath: input.IconPath,
    Name: input.Name,
    AttributeType: input.AttributeType,
    Slug: input.Slug,
}

return output
}
