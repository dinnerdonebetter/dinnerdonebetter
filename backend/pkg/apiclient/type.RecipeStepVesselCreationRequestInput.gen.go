// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepVesselCreationRequestInput struct {
   Name string `json:"name"`
 Notes string `json:"notes"`
 ProductOfRecipeStepIndex uint64 `json:"productOfRecipeStepIndex"`
 ProductOfRecipeStepProductIndex uint64 `json:"productOfRecipeStepProductIndex"`
 Quantity Uint16RangeWithOptionalMax `json:"quantity"`
 RecipeStepProductID string `json:"recipeStepProductID"`
 UnavailableAfterStep bool `json:"unavailableAfterStep"`
 VesselID string `json:"vesselID"`
 VesselPreposition string `json:"vesselPreposition"`

}
)
