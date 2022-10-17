package integration

import (
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanGroceryListItemEquality(t *testing.T, expected, actual *types.MealPlanGroceryListItem) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.QuantityPurchased, actual.QuantityPurchased, "expected QuantityPurchased for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.QuantityPurchased, actual.QuantityPurchased)
	assert.Equal(t, expected.PurchasePrice, actual.PurchasePrice, "expected PurchasePrice for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasePrice, actual.PurchasePrice)
	assert.Equal(t, expected.PurchasedUPC, actual.PurchasedUPC, "expected PurchasedUPC for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasedUPC, actual.PurchasedUPC)
	assert.Equal(t, expected.PurchasedMeasurementUnit, actual.PurchasedMeasurementUnit, "expected PurchasedMeasurementUnit for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.PurchasedMeasurementUnit, actual.PurchasedMeasurementUnit)
	assert.Equal(t, expected.StatusExplanation, actual.StatusExplanation, "expected StatusExplanation for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.StatusExplanation, actual.StatusExplanation)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.Equal(t, expected.MeasurementUnit, actual.MeasurementUnit, "expected MeasurementUnit for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit, actual.MeasurementUnit)
	assert.Equal(t, expected.MealPlanOption, actual.MealPlanOption, "expected MealPlanOption for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MealPlanOption, actual.MealPlanOption)
	assert.Equal(t, expected.Ingredient, actual.Ingredient, "expected Ingredient for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.Ingredient, actual.Ingredient)
	assert.Equal(t, expected.MaximumQuantityNeeded, actual.MaximumQuantityNeeded, "expected MaximumQuantityNeeded for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MaximumQuantityNeeded, actual.MaximumQuantityNeeded)
	assert.Equal(t, expected.MinimumQuantityNeeded, actual.MinimumQuantityNeeded, "expected MinimumQuantityNeeded for meal plan grocery list item %s to be %v, but it was %v", expected.ID, expected.MinimumQuantityNeeded, actual.MinimumQuantityNeeded)

	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestMealPlanGroceryListItems_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)

			t.Log("creating meal plan task")
			exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
			exampleMealPlanGroceryListItemInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemCreationRequestInput(exampleMealPlanGroceryListItem)

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleMealPlanGroceryListItemInput.MealPlanOptionID = createdMealPlan.Events[0].Options[0].ID
			exampleMealPlanGroceryListItemInput.ValidIngredientID = createdValidIngredient.ID
			exampleMealPlanGroceryListItemInput.ValidMeasurementUnitID = exampleValidMeasurementUnit.ID

			createdMealPlanGroceryListItem, err := testClients.admin.CreateMealPlanGroceryListItem(ctx, createdMealPlan.ID, exampleMealPlanGroceryListItemInput)
			require.NoError(t, err)
			t.Logf("meal plan task %q created", createdMealPlanGroceryListItem.ID)
			checkMealPlanGroceryListItemEquality(t, exampleMealPlanGroceryListItem, createdMealPlanGroceryListItem)

			t.Log("fetching changed meal plan task")
			actual, err := testClients.admin.GetMealPlanGroceryListItem(ctx, createdMealPlan.ID, createdMealPlanGroceryListItem.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan task equality
			checkMealPlanGroceryListItemEquality(t, exampleMealPlanGroceryListItem, actual)

			t.Log("fetching changed meal plan task")
			actualList, err := testClients.admin.GetMealPlanGroceryListItemsForMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, actual, err)

			assert.Len(t, actualList, 1)
			checkMealPlanGroceryListItemEquality(t, actualList[0], actual)
		}
	})
}
