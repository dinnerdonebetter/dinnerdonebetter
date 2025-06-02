// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientMeasurementUnitUpdateRequestInput struct {
		Notes                  string                                        `json:"notes"`
		ValidIngredientID      string                                        `json:"validIngredientID"`
		ValidMeasurementUnitID string                                        `json:"validMeasurementUnitID"`
		AllowableQuantity      Float32RangeWithOptionalMaxUpdateRequestInput `json:"allowableQuantity"`
	}
)
