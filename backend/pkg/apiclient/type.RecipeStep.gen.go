// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStep struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipe string `json:"belongsToRecipe"`
 CompletionConditions []RecipeStepCompletionCondition `json:"completionConditions"`
 ConditionExpression string `json:"conditionExpression"`
 CreatedAt string `json:"createdAt"`
 EstimatedTimeInSeconds OptionalUint32Range `json:"estimatedTimeInSeconds"`
 ExplicitInstructions string `json:"explicitInstructions"`
 ID string `json:"id"`
 Index uint64 `json:"index"`
 Ingredients []RecipeStepIngredient `json:"ingredients"`
 Instruments []RecipeStepInstrument `json:"instruments"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Media []RecipeMedia `json:"media"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`
 Preparation ValidPreparation `json:"preparation"`
 Products []RecipeStepProduct `json:"products"`
 StartTimerAutomatically bool `json:"startTimerAutomatically"`
 TemperatureInCelsius OptionalFloat32Range `json:"temperatureInCelsius"`
 Vessels []RecipeStepVessel `json:"vessels"`

}
)
