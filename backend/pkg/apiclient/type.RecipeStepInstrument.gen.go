// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepInstrument struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 Instrument ValidInstrument `json:"instrument"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 OptionIndex uint64 `json:"optionIndex"`
 Optional bool `json:"optional"`
 PreferenceRank uint64 `json:"preferenceRank"`
 Quantity Uint32RangeWithOptionalMax `json:"quantity"`
 RecipeStepProductID string `json:"recipeStepProductID"`

}
)
