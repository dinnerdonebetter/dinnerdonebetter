package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"

	"github.com/stretchr/testify/assert"
)

func checkMealPlanGroceryListItemEquality(t *testing.T, expected, actual *mealplanning.MealPlanGroceryListItem) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.QuantityPurchased, actual.QuantityPurchased, "expected QuantityPurchased for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.QuantityPurchased, actual.QuantityPurchased)
	assert.Equal(t, expected.PurchasePrice, actual.PurchasePrice, "expected PurchasePrice for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasePrice, actual.PurchasePrice)
	assert.Equal(t, expected.PurchasedUPC, actual.PurchasedUPC, "expected PurchasedUPC for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasedUPC, actual.PurchasedUPC)
	assert.Equal(t, expected.PurchasedMeasurementUnit, actual.PurchasedMeasurementUnit, "expected PurchasedMeasurementUnit for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasedMeasurementUnit, actual.PurchasedMeasurementUnit)
	assert.Equal(t, expected.StatusExplanation, actual.StatusExplanation, "expected StatusExplanation for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.StatusExplanation, actual.StatusExplanation)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected MeasurementUnitID for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID)
	assert.Equal(t, expected.BelongsToMealPlan, actual.BelongsToMealPlan, "expected BelongsToMealPlan for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.BelongsToMealPlan, actual.BelongsToMealPlan)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected ValidIngredientID for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.Equal(t, expected.QuantityNeeded, actual.QuantityNeeded, "expected QuantityNeeded for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.QuantityNeeded, actual.QuantityNeeded)

	assert.NotZero(t, actual.CreatedAt)
}

func TestMealPlanGroceryListItems_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()

		// TODO: This test should work by first creating a meal plan, finalizing it, and then manipulating its entries

		t.SkipNow()
	})
}
