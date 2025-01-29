// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipePrepTaskWithinRecipeCreationRequestInput struct {
   BelongsToRecipe string `json:"belongsToRecipe"`
 Description string `json:"description"`
 ExplicitStorageInstructions string `json:"explicitStorageInstructions"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`
 RecipeSteps []RecipePrepTaskStepWithinRecipeCreationRequestInput `json:"recipeSteps"`
 StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
 StorageType string `json:"storageType"`
 TimeBufferBeforeRecipeInSeconds Uint32RangeWithOptionalMax `json:"timeBufferBeforeRecipeInSeconds"`

}
)
