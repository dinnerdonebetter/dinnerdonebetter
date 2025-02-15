// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepIngredient struct {
		ProductOfRecipeID      string                      `json:"productOfRecipeID"`
		IngredientNotes        string                      `json:"ingredientNotes"`
		ID                     string                      `json:"id"`
		CreatedAt              string                      `json:"createdAt"`
		BelongsToRecipeStep    string                      `json:"belongsToRecipeStep"`
		QuantityNotes          string                      `json:"quantityNotes"`
		LastUpdatedAt          string                      `json:"lastUpdatedAt"`
		ArchivedAt             string                      `json:"archivedAt"`
		RecipeStepProductID    string                      `json:"recipeStepProductID"`
		Name                   string                      `json:"name"`
		MeasurementUnit        ValidMeasurementUnit        `json:"measurementUnit"`
		Ingredient             ValidIngredient             `json:"ingredient"`
		Quantity               Float32RangeWithOptionalMax `json:"quantity"`
		ProductPercentageToUse float64                     `json:"productPercentageToUse"`
		OptionIndex            uint64                      `json:"optionIndex"`
		VesselIndex            uint64                      `json:"vesselIndex"`
		ToTaste                bool                        `json:"toTaste"`
		Optional               bool                        `json:"optional"`
	}
)
