package converters

import "github.com/prixfixeco/api_server/pkg/types"

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
