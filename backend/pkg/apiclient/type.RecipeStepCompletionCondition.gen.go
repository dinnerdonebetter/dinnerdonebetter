// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepCompletionCondition struct {
		IngredientState     ValidIngredientState                      `json:"ingredientState"`
		ArchivedAt          string                                    `json:"archivedAt"`
		BelongsToRecipeStep string                                    `json:"belongsToRecipeStep"`
		CreatedAt           string                                    `json:"createdAt"`
		ID                  string                                    `json:"id"`
		LastUpdatedAt       string                                    `json:"lastUpdatedAt"`
		Notes               string                                    `json:"notes"`
		Ingredients         []RecipeStepCompletionConditionIngredient `json:"ingredients"`
		Optional            bool                                      `json:"optional"`
	}
)
