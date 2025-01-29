// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepInstrumentCreationRequestInput struct {
   InstrumentID string `json:"instrumentID"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 OptionIndex uint64 `json:"optionIndex"`
 Optional bool `json:"optional"`
 PreferenceRank uint64 `json:"preferenceRank"`
 ProductOfRecipeStepIndex uint64 `json:"productOfRecipeStepIndex"`
 ProductOfRecipeStepProductIndex uint64 `json:"productOfRecipeStepProductIndex"`
 Quantity Uint32RangeWithOptionalMax `json:"quantity"`
 RecipeStepProductID string `json:"recipeStepProductID"`

}
)
