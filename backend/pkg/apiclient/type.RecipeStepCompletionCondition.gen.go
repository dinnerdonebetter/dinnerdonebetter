// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepCompletionCondition struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 IngredientState ValidIngredientState `json:"ingredientState"`
 Ingredients []RecipeStepCompletionConditionIngredient `json:"ingredients"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`

}
)
