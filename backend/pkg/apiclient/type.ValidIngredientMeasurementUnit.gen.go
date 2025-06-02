// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientMeasurementUnit struct {
		ArchivedAt        string                      `json:"archivedAt"`
		CreatedAt         string                      `json:"createdAt"`
		ID                string                      `json:"id"`
		LastUpdatedAt     string                      `json:"lastUpdatedAt"`
		Notes             string                      `json:"notes"`
		MeasurementUnit   ValidMeasurementUnit        `json:"measurementUnit"`
		Ingredient        ValidIngredient             `json:"ingredient"`
		AllowableQuantity Float32RangeWithOptionalMax `json:"allowableQuantity"`
	}
)
