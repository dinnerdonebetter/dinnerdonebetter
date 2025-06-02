// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanGroceryListItemCreationRequestInput struct {
		BelongsToMealPlan          string                      `json:"belongsToMealPlan"`
		PurchasedMeasurementUnitID string                      `json:"purchasedMeasurementUnitID"`
		PurchasedUpc               string                      `json:"purchasedUPC"`
		Status                     string                      `json:"status"`
		StatusExplanation          string                      `json:"statusExplanation"`
		ValidIngredientID          string                      `json:"validIngredientID"`
		ValidMeasurementUnitID     string                      `json:"validMeasurementUnitID"`
		QuantityNeeded             Float32RangeWithOptionalMax `json:"quantityNeeded"`
		PurchasePrice              float64                     `json:"purchasePrice"`
		QuantityPurchased          float64                     `json:"quantityPurchased"`
	}
)
