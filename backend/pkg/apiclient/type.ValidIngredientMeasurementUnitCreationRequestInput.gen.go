// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientMeasurementUnitCreationRequestInput struct {
		Notes                  string                      `json:"notes"`
		ValidIngredientID      string                      `json:"validIngredientID"`
		ValidMeasurementUnitID string                      `json:"validMeasurementUnitID"`
		AllowableQuantity      Float32RangeWithOptionalMax `json:"allowableQuantity"`
	}
)
