// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepUpdateRequestInput struct {
   BelongsToRecipe string `json:"belongsToRecipe"`
 ConditionExpression string `json:"conditionExpression"`
 EstimatedTimeInSeconds OptionalUint32Range `json:"estimatedTimeInSeconds"`
 ExplicitInstructions string `json:"explicitInstructions"`
 Index uint64 `json:"index"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`
 Preparation ValidPreparation `json:"preparation"`
 StartTimerAutomatically bool `json:"startTimerAutomatically"`
 TemperatureInCelsius OptionalFloat32Range `json:"temperatureInCelsius"`

}
)
