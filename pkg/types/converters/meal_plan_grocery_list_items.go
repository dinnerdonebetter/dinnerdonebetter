package converters

import (
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput builds a MealPlanGroceryListItemDatabaseCreationInput from a MealPlanGroceryListItem.
func ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(input *types.MealPlanGroceryListItem) *types.MealPlanGroceryListItemDatabaseCreationInput {
	x := &types.MealPlanGroceryListItemDatabaseCreationInput{
		ID:                     input.ID,
		MealPlanOptionID:       input.MealPlanOption.ID,
		ValidIngredientID:      input.Ingredient.ID,
		ValidMeasurementUnitID: input.MeasurementUnit.ID,
		MinimumQuantityNeeded:  input.MinimumQuantityNeeded,
		MaximumQuantityNeeded:  input.MaximumQuantityNeeded,
		QuantityPurchased:      input.QuantityPurchased,
		PurchasedUPC:           input.PurchasedUPC,
		PurchasePrice:          input.PurchasePrice,
		StatusExplanation:      input.StatusExplanation,
		Status:                 input.Status,
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
		MealPlanOptionID:       input.MealPlanOption.ID,
		ValidIngredientID:      input.Ingredient.ID,
		ValidMeasurementUnitID: input.MeasurementUnit.ID,
		MinimumQuantityNeeded:  input.MinimumQuantityNeeded,
		MaximumQuantityNeeded:  input.MaximumQuantityNeeded,
	}

	if input.PurchasedMeasurementUnit != nil {
		x.PurchasedMeasurementUnitID = &input.PurchasedMeasurementUnit.ID
	}

	return x
}

func ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemDatabaseCreationInput(input *types.MealPlanGroceryListItemCreationRequestInput) *types.MealPlanGroceryListItemDatabaseCreationInput {
	return &types.MealPlanGroceryListItemDatabaseCreationInput{
		PurchasePrice:              input.PurchasePrice,
		PurchasedUPC:               input.PurchasedUPC,
		PurchasedMeasurementUnitID: input.PurchasedMeasurementUnitID,
		QuantityPurchased:          input.QuantityPurchased,
		Status:                     input.Status,
		StatusExplanation:          input.StatusExplanation,
		ValidMeasurementUnitID:     input.ValidMeasurementUnitID,
		ValidIngredientID:          input.ValidIngredientID,
		MealPlanOptionID:           input.MealPlanOptionID,
		MinimumQuantityNeeded:      input.MinimumQuantityNeeded,
		MaximumQuantityNeeded:      input.MaximumQuantityNeeded,
	}
}

func ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(input *types.MealPlanGroceryListItem) *types.MealPlanGroceryListItemUpdateRequestInput {
	var purchasedMeasurementUnitID *string
	if input.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnitID = &input.PurchasedMeasurementUnit.ID
	}

	return &types.MealPlanGroceryListItemUpdateRequestInput{
		ID:                         input.ID,
		MealPlanOptionID:           &input.MealPlanOption.ID,
		ValidIngredientID:          &input.Ingredient.ID,
		ValidMeasurementUnitID:     &input.MeasurementUnit.ID,
		MinimumQuantityNeeded:      &input.MinimumQuantityNeeded,
		MaximumQuantityNeeded:      &input.MaximumQuantityNeeded,
		QuantityPurchased:          input.QuantityPurchased,
		PurchasedMeasurementUnitID: purchasedMeasurementUnitID,
		PurchasedUPC:               input.PurchasedUPC,
		PurchasePrice:              input.PurchasePrice,
		StatusExplanation:          &input.StatusExplanation,
		Status:                     &input.Status,
	}
}