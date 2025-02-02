// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepVesselUpdateRequestInput struct {
		BelongsToRecipeStep  string                                       `json:"belongsToRecipeStep"`
		Name                 string                                       `json:"name"`
		Notes                string                                       `json:"notes"`
		RecipeStepProductID  string                                       `json:"recipeStepProductID"`
		VesselID             string                                       `json:"vesselID"`
		VesselPreposition    string                                       `json:"vesselPreposition"`
		Quantity             Uint16RangeWithOptionalMaxUpdateRequestInput `json:"quantity"`
		UnavailableAfterStep bool                                         `json:"unavailableAfterStep"`
	}
)
