package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeCreationRequestInputToRecipe(input *messages.RecipeCreationRequestInput) *messages.Recipe {
convertedprepTasks := make([]*messages.RecipePrepTask, 0, len(input.PrepTasks))
for _, item := range input.PrepTasks {
    convertedprepTasks = append(convertedprepTasks, ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTask(item))
}
convertedsteps := make([]*messages.RecipeStep, 0, len(input.Steps))
for _, item := range input.Steps {
    convertedsteps = append(convertedsteps, ConvertRecipeStepCreationRequestInputToRecipeStep(item))
}

output := &messages.Recipe{
    EstimatedPortions: input.EstimatedPortions,
    YieldsComponentType: input.YieldsComponentType,
    Name: input.Name,
    PortionName: input.PortionName,
    Slug: input.Slug,
    PluralPortionName: input.PluralPortionName,
    Description: input.Description,
    InspiredByRecipeID: input.InspiredByRecipeID,
    Source: input.Source,
    PrepTasks: convertedprepTasks,
    Steps: convertedsteps,
    SealOfApproval: input.SealOfApproval,
    EligibleForMeals: input.EligibleForMeals,
}

return output
}
func ConvertRecipeUpdateRequestInputToRecipe(input *messages.RecipeUpdateRequestInput) *messages.Recipe {

output := &messages.Recipe{
    Name: input.Name,
    PortionName: input.PortionName,
    InspiredByRecipeID: input.InspiredByRecipeID,
    Slug: input.Slug,
    PluralPortionName: input.PluralPortionName,
    EstimatedPortions: ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input.EstimatedPortions),
    YieldsComponentType: input.YieldsComponentType,
    Description: input.Description,
    Source: input.Source,
    SealOfApproval: input.SealOfApproval,
    EligibleForMeals: input.EligibleForMeals,
}

return output
}
