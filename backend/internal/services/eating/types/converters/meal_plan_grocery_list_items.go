package converters

import (
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

// ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput builds a MealPlanGroceryListItemDatabaseCreationInput from a MealPlanGroceryListItem.
func ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(input *types.MealPlanGroceryListItem) *types.MealPlanGroceryListItemDatabaseCreationInput {
	x := &types.MealPlanGroceryListItemDatabaseCreationInput{
		ID:                     input.ID,
		BelongsToMealPlan:      input.BelongsToMealPlan,
		ValidIngredientID:      input.Ingredient.ID,
		ValidMeasurementUnitID: input.MeasurementUnit.ID,
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
		QuantityPurchased: input.QuantityPurchased,
		PurchasedUPC:      input.PurchasedUPC,
		PurchasePrice:     input.PurchasePrice,
		StatusExplanation: input.StatusExplanation,
		Status:            input.Status,
	}

	if input.PurchasedMeasurementUnit != nil {
		x.PurchasedMeasurementUnitID = &input.PurchasedMeasurementUnit.ID
	}

	return x
}

// ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput builds a MealPlanGroceryListItemCreationRequestInput from a MealPlanGroceryListItem.
func ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(input *types.MealPlanGroceryListItem) *types.MealPlanGroceryListItemCreationRequestInput {
	x := &types.MealPlanGroceryListItemCreationRequestInput{
		PurchasedUPC:           input.PurchasedUPC,
		PurchasePrice:          input.PurchasePrice,
		QuantityPurchased:      input.QuantityPurchased,
		StatusExplanation:      input.StatusExplanation,
		Status:                 input.Status,
		BelongsToMealPlan:      input.BelongsToMealPlan,
		ValidIngredientID:      input.Ingredient.ID,
		ValidMeasurementUnitID: input.MeasurementUnit.ID,
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
	}

	if input.PurchasedMeasurementUnit != nil {
		x.PurchasedMeasurementUnitID = &input.PurchasedMeasurementUnit.ID
	}

	return x
}

func ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemDatabaseCreationInput(input *types.MealPlanGroceryListItemCreationRequestInput) *types.MealPlanGroceryListItemDatabaseCreationInput {
	return &types.MealPlanGroceryListItemDatabaseCreationInput{
		ID:                         identifiers.New(),
		PurchasePrice:              input.PurchasePrice,
		PurchasedUPC:               input.PurchasedUPC,
		PurchasedMeasurementUnitID: input.PurchasedMeasurementUnitID,
		QuantityPurchased:          input.QuantityPurchased,
		Status:                     input.Status,
		StatusExplanation:          input.StatusExplanation,
		ValidMeasurementUnitID:     input.ValidMeasurementUnitID,
		ValidIngredientID:          input.ValidIngredientID,
		BelongsToMealPlan:          input.BelongsToMealPlan,
		QuantityNeeded: types.Float32RangeWithOptionalMax{
			Max: input.QuantityNeeded.Max,
			Min: input.QuantityNeeded.Min,
		},
	}
}

func ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(input *types.MealPlanGroceryListItem) *types.MealPlanGroceryListItemUpdateRequestInput {
	var purchasedMeasurementUnitID *string
	if input.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnitID = &input.PurchasedMeasurementUnit.ID
	}

	return &types.MealPlanGroceryListItemUpdateRequestInput{
		BelongsToMealPlan:      &input.BelongsToMealPlan,
		ValidIngredientID:      &input.Ingredient.ID,
		ValidMeasurementUnitID: &input.MeasurementUnit.ID,
		QuantityNeeded: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Max: input.QuantityNeeded.Max,
			Min: &input.QuantityNeeded.Min,
		},
		QuantityPurchased:          input.QuantityPurchased,
		PurchasedMeasurementUnitID: purchasedMeasurementUnitID,
		PurchasedUPC:               input.PurchasedUPC,
		PurchasePrice:              input.PurchasePrice,
		StatusExplanation:          &input.StatusExplanation,
		Status:                     &input.Status,
	}
}
