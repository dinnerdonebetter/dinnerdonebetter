package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipePrepTaskCreationRequestInputToRecipePrepTask(input *messages.RecipePrepTaskCreationRequestInput) *messages.RecipePrepTask {

output := &messages.RecipePrepTask{
    StorageType: input.StorageType,
    ExplicitStorageInstructions: input.ExplicitStorageInstructions,
    Notes: input.Notes,
    Name: input.Name,
    Description: input.Description,
    Optional: input.Optional,
    StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
    TimeBufferBeforeRecipeInSeconds: input.TimeBufferBeforeRecipeInSeconds,
    BelongsToRecipe: input.BelongsToRecipe,
}

return output
}
func ConvertRecipePrepTaskUpdateRequestInputToRecipePrepTask(input *messages.RecipePrepTaskUpdateRequestInput) *messages.RecipePrepTask {
convertedtaskSteps := make([]*messages.RecipePrepTaskStep, 0, len(input.TaskSteps))
for _, item := range input.TaskSteps {
    convertedtaskSteps = append(convertedtaskSteps, ConvertRecipePrepTaskStepUpdateRequestInputToRecipePrepTaskStep(item))
}

output := &messages.RecipePrepTask{
    StorageType: input.StorageType,
    TaskSteps: convertedtaskSteps,
    Optional: input.Optional,
    StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
    TimeBufferBeforeRecipeInSeconds: ConvertUint32RangeWithOptionalMaxUpdateRequestInputToUint32RangeWithOptionalMax(input.TimeBufferBeforeRecipeInSeconds),
    BelongsToRecipe: input.BelongsToRecipe,
    ExplicitStorageInstructions: input.ExplicitStorageInstructions,
    Notes: input.Notes,
    Name: input.Name,
    Description: input.Description,
}

return output
}
