// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepCompletionConditionForExistingRecipeCreationRequestInput struct {
   BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 IngredientStateID string `json:"ingredientStateID"`
 Ingredients []RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput `json:"ingredients"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`

}
)
