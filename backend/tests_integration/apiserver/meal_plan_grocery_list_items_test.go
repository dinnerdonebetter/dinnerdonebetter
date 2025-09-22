package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		ctx := t.Context()

		createdMealPlan := createMealPlanForTest(t, adminClient, nil)

		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleMealPlanGroceryListItemInput := mpconverters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(exampleMealPlanGroceryListItem)

		createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)
		createdValidIngredient := createValidIngredientForTest(t)

		exampleMealPlanGroceryListItem.BelongsToMealPlan = createdMealPlan.ID
		exampleMealPlanGroceryListItem.Ingredient = *createdValidIngredient
		exampleMealPlanGroceryListItem.MeasurementUnit = *createdValidMeasurementUnit

		exampleMealPlanGroceryListItemInput.BelongsToMealPlan = createdMealPlan.ID
		exampleMealPlanGroceryListItemInput.ValidIngredientID = createdValidIngredient.ID
		exampleMealPlanGroceryListItemInput.ValidMeasurementUnitID = createdValidMeasurementUnit.ID

		createdMealPlanGroceryListItemRes, err := adminClient.CreateMealPlanGroceryListItem(ctx, &mealplanninggrpc.CreateMealPlanGroceryListItemRequest{
			MealPlanID: createdMealPlan.ID,
			Input:      converters.ConvertMealPlanGroceryListItemCreationRequestInputToGRPCMealPlanGroceryListItemCreationRequestInput(exampleMealPlanGroceryListItemInput),
		})
		require.NoError(t, err)

		createdMealPlanGroceryListItem := converters.ConvertGRPCMealPlanGroceryListItemToMealPlanGroceryListItem(createdMealPlanGroceryListItemRes.Created)
		checkMealPlanGroceryListItemEquality(t, exampleMealPlanGroceryListItem, createdMealPlanGroceryListItem)

		actualRes, err := adminClient.GetMealPlanGroceryListItem(ctx, &mealplanninggrpc.GetMealPlanGroceryListItemRequest{
			MealPlanID:                createdMealPlan.ID,
			MealPlanGroceryListItemID: createdMealPlanGroceryListItem.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, actualRes)

		// assert meal plan task equality
		actual := converters.ConvertGRPCMealPlanGroceryListItemToMealPlanGroceryListItem(actualRes.Result)
		checkMealPlanGroceryListItemEquality(t, exampleMealPlanGroceryListItem, actual)

		actualList, err := adminClient.GetMealPlanGroceryListItemsForMealPlan(ctx, &mealplanninggrpc.GetMealPlanGroceryListItemsForMealPlanRequest{
			Filter:     nil,
			MealPlanID: createdMealPlan.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, actualList)

		_, err = adminClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanID: createdMealPlan.ID})
		assert.NoError(t, err)

		assert.Len(t, actualList.Results, 1)
		checkMealPlanGroceryListItemEquality(t, converters.ConvertGRPCMealPlanGroceryListItemToMealPlanGroceryListItem(actualList.Results[0]), actual)
	})
}
