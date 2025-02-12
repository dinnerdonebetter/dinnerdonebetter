package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItem(input *messages.MealPlanGroceryListItemCreationRequestInput) *messages.MealPlanGroceryListItem {

output := &messages.MealPlanGroceryListItem{
    QuantityPurchased: input.QuantityPurchased,
    QuantityNeeded: input.QuantityNeeded,
    PurchasedUPC: input.PurchasedUPC,
    Status: input.Status,
    StatusExplanation: input.StatusExplanation,
    BelongsToMealPlan: input.BelongsToMealPlan,
    PurchasePrice: input.PurchasePrice,
}

return output
}
func ConvertMealPlanGroceryListItemUpdateRequestInputToMealPlanGroceryListItem(input *messages.MealPlanGroceryListItemUpdateRequestInput) *messages.MealPlanGroceryListItem {

output := &messages.MealPlanGroceryListItem{
    PurchasePrice: input.PurchasePrice,
    QuantityPurchased: input.QuantityPurchased,
    QuantityNeeded: ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input.QuantityNeeded),
    PurchasedUPC: input.PurchasedUPC,
    Status: input.Status,
    StatusExplanation: input.StatusExplanation,
    BelongsToMealPlan: input.BelongsToMealPlan,
}

return output
}
