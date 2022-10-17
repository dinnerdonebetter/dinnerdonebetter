package fakes

import (
	"time"

	"github.com/prixfixeco/api_server/pkg/types"
)

func BuildFakeMealPlanGroceryListItem() *types.MealPlanGroceryListItem {
	minQty := BuildFakeNumber()

	return &types.MealPlanGroceryListItem{
		ID:                       BuildFakeID(),
		MealPlanOption:           *BuildFakeMealPlanOption(),
		Ingredient:               *BuildFakeValidIngredient(),
		MeasurementUnit:          *BuildFakeValidMeasurementUnit(),
		MinimumQuantityNeeded:    float32(minQty),
		MaximumQuantityNeeded:    float32(minQty + 1),
		QuantityPurchased:        nil,
		PurchasedMeasurementUnit: nil,
		PurchasedUPC:             nil,
		PurchasePrice:            nil,
		StatusExplanation:        buildUniqueString(),
		Status:                   "unknown",
		CreatedAt:                time.Now(),
	}
}

func BuildFakeMealPlanGroceryListItemList() *types.MealPlanGroceryListItemList {
	var examples []*types.MealPlanGroceryListItem
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanGroceryListItem())
	}

	return &types.MealPlanGroceryListItemList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlanGroceryListItems: examples,
	}
}
