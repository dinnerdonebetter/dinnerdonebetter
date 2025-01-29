// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepInstrumentUpdateRequestInput struct {
   BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 InstrumentID string `json:"instrumentID"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 OptionIndex uint64 `json:"optionIndex"`
 Optional bool `json:"optional"`
 PreferenceRank uint64 `json:"preferenceRank"`
 Quantity Uint32RangeWithOptionalMaxUpdateRequestInput `json:"quantity"`
 RecipeStepProductID string `json:"recipeStepProductID"`

}
)
