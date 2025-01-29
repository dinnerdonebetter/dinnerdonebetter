// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeStepVessel struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToRecipeStep string `json:"belongsToRecipeStep"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Name string `json:"name"`
 Notes string `json:"notes"`
 Quantity Uint16RangeWithOptionalMax `json:"quantity"`
 RecipeStepProductID string `json:"recipeStepProductID"`
 UnavailableAfterStep bool `json:"unavailableAfterStep"`
 Vessel ValidVessel `json:"vessel"`
 VesselPreposition string `json:"vesselPreposition"`

}
)
