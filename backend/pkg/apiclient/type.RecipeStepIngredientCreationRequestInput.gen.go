// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepIngredientCreationRequestInput struct {
   IngredientID string `json:"ingredientID"`
 IngredientNotes string `json:"ingredientNotes"`
 MeasurementUnitID string `json:"measurementUnitID"`
 Name string `json:"name"`
 OptionIndex uint64 `json:"optionIndex"`
 Optional bool `json:"optional"`
 ProductOfRecipeID string `json:"productOfRecipeID"`
 ProductOfRecipeStepIndex uint64 `json:"productOfRecipeStepIndex"`
 ProductOfRecipeStepProductIndex uint64 `json:"productOfRecipeStepProductIndex"`
 ProductPercentageToUse float64 `json:"productPercentageToUse"`
 Quantity Float32RangeWithOptionalMax `json:"quantity"`
 QuantityNotes string `json:"quantityNotes"`
 ToTaste bool `json:"toTaste"`
 VesselIndex uint64 `json:"vesselIndex"`

}
)
