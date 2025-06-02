// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanGroceryListItem struct {
		ArchivedAt               string                      `json:"archivedAt"`
		PurchasedUpc             string                      `json:"purchasedUPC"`
		CreatedAt                string                      `json:"createdAt"`
		ID                       string                      `json:"id"`
		BelongsToMealPlan        string                      `json:"belongsToMealPlan"`
		LastUpdatedAt            string                      `json:"lastUpdatedAt"`
		StatusExplanation        string                      `json:"statusExplanation"`
		Status                   string                      `json:"status"`
		PurchasedMeasurementUnit ValidMeasurementUnit        `json:"purchasedMeasurementUnit"`
		MeasurementUnit          ValidMeasurementUnit        `json:"measurementUnit"`
		Ingredient               ValidIngredient             `json:"ingredient"`
		QuantityNeeded           Float32RangeWithOptionalMax `json:"quantityNeeded"`
		PurchasePrice            float64                     `json:"purchasePrice"`
		QuantityPurchased        float64                     `json:"quantityPurchased"`
	}
)
