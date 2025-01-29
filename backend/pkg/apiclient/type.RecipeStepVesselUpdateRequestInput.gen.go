// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepVesselUpdateRequestInput struct {
   BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 Quantity Uint16RangeWithOptionalMaxUpdateRequestInput `json:"quantity"`
 RecipeStepProductID string `json:"recipeStepProductID"`
 UnavailableAfterStep bool `json:"unavailableAfterStep"`
 VesselID string `json:"vesselID"`
 VesselPreposition string `json:"vesselPreposition"`

}
)
