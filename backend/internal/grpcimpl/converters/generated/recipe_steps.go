package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepCreationRequestInputToRecipeStep(input *messages.RecipeStepCreationRequestInput) *messages.RecipeStep {
convertedinstruments := make([]*messages.RecipeStepInstrument, 0, len(input.Instruments))
for _, item := range input.Instruments {
    convertedinstruments = append(convertedinstruments, ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrument(item))
}
convertedvessels := make([]*messages.RecipeStepVessel, 0, len(input.Vessels))
for _, item := range input.Vessels {
    convertedvessels = append(convertedvessels, ConvertRecipeStepVesselCreationRequestInputToRecipeStepVessel(item))
}
convertedcompletionConditions := make([]*messages.RecipeStepCompletionCondition, 0, len(input.CompletionConditions))
for _, item := range input.CompletionConditions {
    convertedcompletionConditions = append(convertedcompletionConditions, ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionCondition(item))
}
convertedingredients := make([]*messages.RecipeStepIngredient, 0, len(input.Ingredients))
for _, item := range input.Ingredients {
    convertedingredients = append(convertedingredients, ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredient(item))
}
convertedproducts := make([]*messages.RecipeStepProduct, 0, len(input.Products))
for _, item := range input.Products {
    convertedproducts = append(convertedproducts, ConvertRecipeStepProductCreationRequestInputToRecipeStepProduct(item))
}

output := &messages.RecipeStep{
    TemperatureInCelsius: input.TemperatureInCelsius,
    Notes: input.Notes,
    Instruments: convertedinstruments,
    Vessels: convertedvessels,
    CompletionConditions: convertedcompletionConditions,
    Optional: input.Optional,
    EstimatedTimeInSeconds: input.EstimatedTimeInSeconds,
    ExplicitInstructions: input.ExplicitInstructions,
    ConditionExpression: input.ConditionExpression,
    Ingredients: convertedingredients,
    Products: convertedproducts,
    Index: input.Index,
    StartTimerAutomatically: input.StartTimerAutomatically,
}

return output
}
func ConvertRecipeStepUpdateRequestInputToRecipeStep(input *messages.RecipeStepUpdateRequestInput) *messages.RecipeStep {

output := &messages.RecipeStep{
    EstimatedTimeInSeconds: input.EstimatedTimeInSeconds,
    TemperatureInCelsius: input.TemperatureInCelsius,
    Preparation: input.Preparation,
    ExplicitInstructions: input.ExplicitInstructions,
    Notes: input.Notes,
    Optional: input.Optional,
    BelongsToRecipe: input.BelongsToRecipe,
    ConditionExpression: input.ConditionExpression,
    Index: input.Index,
    StartTimerAutomatically: input.StartTimerAutomatically,
}

return output
}
