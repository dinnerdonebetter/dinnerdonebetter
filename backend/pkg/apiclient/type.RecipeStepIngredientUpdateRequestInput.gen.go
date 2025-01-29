// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepIngredientUpdateRequestInput struct {
   BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 IngredientID string `json:"ingredientID"`
 IngredientNotes string `json:"ingredientNotes"`
 MeasurementUnitID string `json:"measurementUnitID"`
 Name string `json:"name"`
 OptionIndex uint64 `json:"optionIndex"`
 Optional bool `json:"optional"`
 ProductOfRecipeID string `json:"productOfRecipeID"`
 ProductPercentageToUse float64 `json:"productPercentageToUse"`
 Quantity Float32RangeWithOptionalMaxUpdateRequestInput `json:"quantity"`
 QuantityNotes string `json:"quantityNotes"`
 RecipeStepProductID string `json:"recipeStepProductID"`
 ToTaste bool `json:"toTaste"`
 VesselIndex uint64 `json:"vesselIndex"`

}
)
