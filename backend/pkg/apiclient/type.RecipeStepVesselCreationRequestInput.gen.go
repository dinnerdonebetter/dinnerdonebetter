// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepVesselCreationRequestInput struct {
		Name                            string                     `json:"name"`
		Notes                           string                     `json:"notes"`
		RecipeStepProductID             string                     `json:"recipeStepProductID"`
		VesselID                        string                     `json:"vesselID"`
		VesselPreposition               string                     `json:"vesselPreposition"`
		Quantity                        Uint16RangeWithOptionalMax `json:"quantity"`
		ProductOfRecipeStepIndex        uint64                     `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex uint64                     `json:"productOfRecipeStepProductIndex"`
		UnavailableAfterStep            bool                       `json:"unavailableAfterStep"`
	}
)
