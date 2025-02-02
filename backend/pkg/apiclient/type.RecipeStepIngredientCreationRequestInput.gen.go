// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepIngredientCreationRequestInput struct {
		ProductOfRecipeID               string                      `json:"productOfRecipeID"`
		IngredientID                    string                      `json:"ingredientID"`
		MeasurementUnitID               string                      `json:"measurementUnitID"`
		Name                            string                      `json:"name"`
		QuantityNotes                   string                      `json:"quantityNotes"`
		IngredientNotes                 string                      `json:"ingredientNotes"`
		Quantity                        Float32RangeWithOptionalMax `json:"quantity"`
		OptionIndex                     uint64                      `json:"optionIndex"`
		ProductPercentageToUse          float64                     `json:"productPercentageToUse"`
		ProductOfRecipeStepProductIndex uint64                      `json:"productOfRecipeStepProductIndex"`
		ProductOfRecipeStepIndex        uint64                      `json:"productOfRecipeStepIndex"`
		VesselIndex                     uint64                      `json:"vesselIndex"`
		Optional                        bool                        `json:"optional"`
		ToTaste                         bool                        `json:"toTaste"`
	}
)
