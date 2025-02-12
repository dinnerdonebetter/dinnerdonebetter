package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepVesselCreationRequestInputToRecipeStepVessel(input *messages.RecipeStepVesselCreationRequestInput) *messages.RecipeStepVessel {

output := &messages.RecipeStepVessel{
    Quantity: input.Quantity,
    Notes: input.Notes,
    RecipeStepProductID: input.RecipeStepProductID,
    VesselPreposition: input.VesselPreposition,
    Name: input.Name,
    UnavailableAfterStep: input.UnavailableAfterStep,
}

return output
}
func ConvertRecipeStepVesselUpdateRequestInputToRecipeStepVessel(input *messages.RecipeStepVesselUpdateRequestInput) *messages.RecipeStepVessel {

output := &messages.RecipeStepVessel{
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    VesselPreposition: input.VesselPreposition,
    Name: input.Name,
    UnavailableAfterStep: input.UnavailableAfterStep,
    Quantity: ConvertUint16RangeWithOptionalMaxUpdateRequestInputToUint16RangeWithOptionalMax(input.Quantity),
    Notes: input.Notes,
    RecipeStepProductID: input.RecipeStepProductID,
}

return output
}
