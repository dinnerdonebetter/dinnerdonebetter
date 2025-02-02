// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepCompletionConditionForExistingRecipeCreationRequestInput struct {
		BelongsToRecipeStep string                                                                         `json:"belongsToRecipeStep"`
		IngredientStateID   string                                                                         `json:"ingredientStateID"`
		Notes               string                                                                         `json:"notes"`
		Ingredients         []RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput `json:"ingredients"`
		Optional            bool                                                                           `json:"optional"`
	}
)
