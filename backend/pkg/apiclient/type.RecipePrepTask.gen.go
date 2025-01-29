// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipePrepTask struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipe string `json:"belongsToRecipe"`
 CreatedAt string `json:"createdAt"`
 Description string `json:"description"`
 ExplicitStorageInstructions string `json:"explicitStorageInstructions"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`
 RecipeSteps []RecipePrepTaskStep `json:"recipeSteps"`
 StorageTemperatureInCelsius OptionalFloat32Range `json:"storageTemperatureInCelsius"`
 StorageType string `json:"storageType"`
 TimeBufferBeforeRecipeInSeconds Uint32RangeWithOptionalMax `json:"timeBufferBeforeRecipeInSeconds"`

}
)
