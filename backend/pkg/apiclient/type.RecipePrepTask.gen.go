// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipePrepTask struct {
		LastUpdatedAt                   string                     `json:"lastUpdatedAt"`
		ArchivedAt                      string                     `json:"archivedAt"`
		CreatedAt                       string                     `json:"createdAt"`
		Description                     string                     `json:"description"`
		ExplicitStorageInstructions     string                     `json:"explicitStorageInstructions"`
		ID                              string                     `json:"id"`
		BelongsToRecipe                 string                     `json:"belongsToRecipe"`
		Name                            string                     `json:"name"`
		StorageType                     string                     `json:"storageType"`
		Notes                           string                     `json:"notes"`
		RecipeSteps                     []RecipePrepTaskStep       `json:"recipeSteps"`
		StorageTemperatureInCelsius     OptionalFloat32Range       `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds Uint32RangeWithOptionalMax `json:"timeBufferBeforeRecipeInSeconds"`
		Optional                        bool                       `json:"optional"`
	}
)
