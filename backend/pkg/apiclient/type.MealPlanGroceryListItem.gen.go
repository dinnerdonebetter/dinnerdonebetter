// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealPlanGroceryListItem struct {
   ArchivedAt string `json:"archivedAt"`
 BelongsToMealPlan string `json:"belongsToMealPlan"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 Ingredient ValidIngredient `json:"ingredient"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 MeasurementUnit ValidMeasurementUnit `json:"measurementUnit"`
 PurchasePrice float64 `json:"purchasePrice"`
 PurchasedMeasurementUnit ValidMeasurementUnit `json:"purchasedMeasurementUnit"`
 PurchasedUpc string `json:"purchasedUPC"`
 QuantityNeeded Float32RangeWithOptionalMax `json:"quantityNeeded"`
 QuantityPurchased float64 `json:"quantityPurchased"`
 Status string `json:"status"`
 StatusExplanation string `json:"statusExplanation"`

}
)
