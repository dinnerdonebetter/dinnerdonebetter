// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanGroceryListItemID := fake.BuildFakeID()

		data := &MealPlanGroceryListItem{}
		expected := &APIResponse[*MealPlanGroceryListItem]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanGroceryListItemID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanGroceryListItemID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanGroceryListItem(ctx, "", mealPlanGroceryListItemID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanGroceryListItem ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanGroceryListItem(ctx, mealPlanID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanGroceryListItemID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanGroceryListItemID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanGroceryListItemID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
