package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionCondition(input *messages.RecipeStepCompletionConditionCreationRequestInput) *messages.RecipeStepCompletionCondition {
convertedingredients := make([]*messages.RecipeStepCompletionConditionIngredient, 0, len(input.Ingredients))
for _, item := range input.Ingredients {
    convertedingredients = append(convertedingredients, Convertuint64ToRecipeStepCompletionConditionIngredient(item))
}

output := &messages.RecipeStepCompletionCondition{
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    Notes: input.Notes,
    Ingredients: convertedingredients,
    Optional: input.Optional,
}

return output
}
func ConvertRecipeStepCompletionConditionUpdateRequestInputToRecipeStepCompletionCondition(input *messages.RecipeStepCompletionConditionUpdateRequestInput) *messages.RecipeStepCompletionCondition {

output := &messages.RecipeStepCompletionCondition{
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    Notes: input.Notes,
    Optional: input.Optional,
}

return output
}
