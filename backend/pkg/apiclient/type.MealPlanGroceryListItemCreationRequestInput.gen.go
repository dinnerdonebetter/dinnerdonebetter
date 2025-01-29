// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealPlanGroceryListItemCreationRequestInput struct {
   BelongsToMealPlan string `json:"belongsToMealPlan"`
 PurchasePrice float64 `json:"purchasePrice"`
 PurchasedMeasurementUnitID string `json:"purchasedMeasurementUnitID"`
 PurchasedUpc string `json:"purchasedUPC"`
 QuantityNeeded Float32RangeWithOptionalMax `json:"quantityNeeded"`
 QuantityPurchased float64 `json:"quantityPurchased"`
 Status string `json:"status"`
 StatusExplanation string `json:"statusExplanation"`
 ValidIngredientID string `json:"validIngredientID"`
 ValidMeasurementUnitID string `json:"validMeasurementUnitID"`

}
)
