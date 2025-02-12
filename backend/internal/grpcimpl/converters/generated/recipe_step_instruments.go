package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrument(input *messages.RecipeStepInstrumentCreationRequestInput) *messages.RecipeStepInstrument {

output := &messages.RecipeStepInstrument{
    RecipeStepProductID: input.RecipeStepProductID,
    OptionIndex: input.OptionIndex,
    PreferenceRank: input.PreferenceRank,
    Optional: input.Optional,
    Quantity: input.Quantity,
    Name: input.Name,
    Notes: input.Notes,
}

return output
}
func ConvertRecipeStepInstrumentUpdateRequestInputToRecipeStepInstrument(input *messages.RecipeStepInstrumentUpdateRequestInput) *messages.RecipeStepInstrument {

output := &messages.RecipeStepInstrument{
    Name: input.Name,
    Notes: input.Notes,
    RecipeStepProductID: input.RecipeStepProductID,
    OptionIndex: input.OptionIndex,
    PreferenceRank: input.PreferenceRank,
    Optional: input.Optional,
    Quantity: ConvertUint32RangeWithOptionalMaxUpdateRequestInputToUint32RangeWithOptionalMax(input.Quantity),
    BelongsToRecipeStep: input.BelongsToRecipeStep,
}

return output
}
