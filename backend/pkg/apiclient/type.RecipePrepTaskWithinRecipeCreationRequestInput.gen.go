// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipePrepTaskWithinRecipeCreationRequestInput struct {
		BelongsToRecipe                 string                                               `json:"belongsToRecipe"`
		Description                     string                                               `json:"description"`
		ExplicitStorageInstructions     string                                               `json:"explicitStorageInstructions"`
		Name                            string                                               `json:"name"`
		Notes                           string                                               `json:"notes"`
		StorageType                     string                                               `json:"storageType"`
		RecipeSteps                     []RecipePrepTaskStepWithinRecipeCreationRequestInput `json:"recipeSteps"`
		StorageTemperatureInCelsius     OptionalFloat32Range                                 `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds Uint32RangeWithOptionalMax                           `json:"timeBufferBeforeRecipeInSeconds"`
		Optional                        bool                                                 `json:"optional"`
	}
)
