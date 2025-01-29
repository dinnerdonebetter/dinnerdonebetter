// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepCompletionConditionCreationRequestInput struct {
   BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 IngredientState string `json:"ingredientState"`
 Ingredients []uint64 `json:"ingredients"`
 Notes string `json:"notes"`
 Optional bool `json:"optional"`

}
)
