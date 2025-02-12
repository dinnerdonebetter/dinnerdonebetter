package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipePrepTaskStepCreationRequestInputToRecipePrepTaskStep(input *messages.RecipePrepTaskStepCreationRequestInput) *messages.RecipePrepTaskStep {

output := &messages.RecipePrepTaskStep{
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    SatisfiesRecipeStep: input.SatisfiesRecipeStep,
}

return output
}
func ConvertRecipePrepTaskStepUpdateRequestInputToRecipePrepTaskStep(input *messages.RecipePrepTaskStepUpdateRequestInput) *messages.RecipePrepTaskStep {

output := &messages.RecipePrepTaskStep{
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
    SatisfiesRecipeStep: input.SatisfiesRecipeStep,
}

return output
}
