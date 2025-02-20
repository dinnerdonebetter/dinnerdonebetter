// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepIngredient struct {
		ArchivedAt             string                      `json:"archivedAt"`
		ProductOfRecipeID      string                      `json:"productOfRecipeID"`
		QuantityNotes          string                      `json:"quantityNotes"`
		ID                     string                      `json:"id"`
		BelongsToRecipeStep    string                      `json:"belongsToRecipeStep"`
		IngredientNotes        string                      `json:"ingredientNotes"`
		LastUpdatedAt          string                      `json:"lastUpdatedAt"`
		RecipeStepProductID    string                      `json:"recipeStepProductID"`
		Name                   string                      `json:"name"`
		CreatedAt              string                      `json:"createdAt"`
		MeasurementUnit        ValidMeasurementUnit        `json:"measurementUnit"`
		Ingredient             ValidIngredient             `json:"ingredient"`
		Quantity               Float32RangeWithOptionalMax `json:"quantity"`
		OptionIndex            uint64                      `json:"optionIndex"`
		ProductPercentageToUse float64                     `json:"productPercentageToUse"`
		VesselIndex            uint64                      `json:"vesselIndex"`
		Optional               bool                        `json:"optional"`
		ToTaste                bool                        `json:"toTaste"`
	}
)
