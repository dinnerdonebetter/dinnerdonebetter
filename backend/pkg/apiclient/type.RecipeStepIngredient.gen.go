// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepIngredient struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 Ingredient ValidIngredient `json:"ingredient"`
 IngredientNotes string `json:"ingredientNotes"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 MeasurementUnit ValidMeasurementUnit `json:"measurementUnit"`
 Name string `json:"name"`
 OptionIndex uint64 `json:"optionIndex"`
 Optional bool `json:"optional"`
 ProductOfRecipeID string `json:"productOfRecipeID"`
 ProductPercentageToUse float64 `json:"productPercentageToUse"`
 Quantity Float32RangeWithOptionalMax `json:"quantity"`
 QuantityNotes string `json:"quantityNotes"`
 RecipeStepProductID string `json:"recipeStepProductID"`
 ToTaste bool `json:"toTaste"`
 VesselIndex uint64 `json:"vesselIndex"`

}
)
