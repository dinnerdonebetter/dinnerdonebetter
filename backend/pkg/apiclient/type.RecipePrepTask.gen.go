// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipePrepTask struct {
		LastUpdatedAt                   string                     `json:"lastUpdatedAt"`
		StorageType                     string                     `json:"storageType"`
		CreatedAt                       string                     `json:"createdAt"`
		Description                     string                     `json:"description"`
		ExplicitStorageInstructions     string                     `json:"explicitStorageInstructions"`
		ID                              string                     `json:"id"`
		Notes                           string                     `json:"notes"`
		ArchivedAt                      string                     `json:"archivedAt"`
		BelongsToRecipe                 string                     `json:"belongsToRecipe"`
		Name                            string                     `json:"name"`
		RecipeSteps                     []RecipePrepTaskStep       `json:"recipeSteps"`
		StorageTemperatureInCelsius     OptionalFloat32Range       `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds Uint32RangeWithOptionalMax `json:"timeBufferBeforeRecipeInSeconds"`
		Optional                        bool                       `json:"optional"`
	}
)
