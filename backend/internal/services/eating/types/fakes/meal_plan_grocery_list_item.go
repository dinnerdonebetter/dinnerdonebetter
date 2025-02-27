package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

func BuildFakeMealPlanGroceryListItem() *types.MealPlanGroceryListItem {
	return &types.MealPlanGroceryListItem{
		ID:                       BuildFakeID(),
		BelongsToMealPlan:        BuildFakeID(),
		Ingredient:               *BuildFakeValidIngredient(),
		MeasurementUnit:          *BuildFakeValidMeasurementUnit(),
		QuantityNeeded:           BuildFakeFloat32RangeWithOptionalMax(),
		QuantityPurchased:        nil,
		PurchasedMeasurementUnit: nil,
		PurchasedUPC:             nil,
		PurchasePrice:            nil,
		StatusExplanation:        buildUniqueString(),
		Status:                   types.MealPlanGroceryListItemStatusUnknown,
		CreatedAt:                BuildFakeTime(),
	}
}

func BuildFakeMealPlanGroceryListItemsList() *filtering.QueryFilteredResult[types.MealPlanGroceryListItem] {
	var examples []*types.MealPlanGroceryListItem
	for range exampleQuantity {
		examples = append(examples, BuildFakeMealPlanGroceryListItem())
	}

	return &filtering.QueryFilteredResult[types.MealPlanGroceryListItem]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

func BuildFakeMealPlanGroceryListItemCreationRequestInput() *types.MealPlanGroceryListItemCreationRequestInput {
	mealPlanGroceryListItem := BuildFakeMealPlanGroceryListItem()
	return converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(mealPlanGroceryListItem)
}

func BuildFakeMealPlanGroceryListItemUpdateRequestInput() *types.MealPlanGroceryListItemUpdateRequestInput {
	mealPlanGroceryListItem := BuildFakeMealPlanGroceryListItem()
	return converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(mealPlanGroceryListItem)
}
