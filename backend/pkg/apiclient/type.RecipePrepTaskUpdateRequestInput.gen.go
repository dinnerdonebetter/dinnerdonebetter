// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipePrepTaskUpdateRequestInput struct {
   BelongsToRecipe string `json:"belongsToRecipe"`
 Description string `json:"description"`
 ExplicitStorageInstructions string `json:"explicitStorageInstructions"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`
 RecipeSteps []RecipePrepTaskStepUpdateRequestInput `json:"recipeSteps"`
 StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
 StorageType string `json:"storageType"`
 TimeBufferBeforeRecipeInSeconds Uint32RangeWithOptionalMaxUpdateRequestInput `json:"timeBufferBeforeRecipeInSeconds"`

}
)
