package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionCondition(input *messages.RecipeStepCompletionConditionCreationRequestInput) *messages.RecipeStepCompletionCondition {
	convertedIngredients := make([]*messages.RecipeStepCompletionConditionIngredient, 0, len(input.Ingredients))
	for _, item := range input.Ingredients {
		convertedIngredients = append(convertedIngredients, &messages.RecipeStepCompletionConditionIngredient{
			ID: item,
		})
	}

	output := &messages.RecipeStepCompletionCondition{
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Ingredients:         convertedIngredients,
		Optional:            input.Optional,
	}

	return output
}

func ConvertRecipeStepCompletionConditionUpdateRequestInputToRecipeStepCompletionCondition(input *messages.RecipeStepCompletionConditionUpdateRequestInput) *messages.RecipeStepCompletionCondition {

	output := &messages.RecipeStepCompletionCondition{
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Notes:               input.Notes,
		Optional:            input.Optional,
	}

	return output
}
