// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepCreationRequestInput struct {
   CompletionConditions []RecipeStepCompletionConditionCreationRequestInput `json:"completionConditions"`
 ConditionExpression string `json:"conditionExpression"`
 EstimatedTimeInSeconds OptionalUint32Range `json:"estimatedTimeInSeconds"`
 ExplicitInstructions string `json:"explicitInstructions"`
 Index uint64 `json:"index"`
 Ingredients []RecipeStepIngredientCreationRequestInput `json:"ingredients"`
 Instruments []RecipeStepInstrumentCreationRequestInput `json:"instruments"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`
 PreparationID string `json:"preparationID"`
 Products []RecipeStepProductCreationRequestInput `json:"products"`
 StartTimerAutomatically bool `json:"startTimerAutomatically"`
 TemperatureInCelsius OptionalFloat32Range `json:"temperatureInCelsius"`
 Vessels []RecipeStepVesselCreationRequestInput `json:"vessels"`

}
)
