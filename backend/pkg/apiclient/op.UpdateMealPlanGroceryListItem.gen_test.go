// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanGroceryListItemID := fake.BuildFakeID()

		data := &MealPlanGroceryListItem{}
		expected := &APIResponse[*MealPlanGroceryListItem]{
			Data: data,
		}

		exampleInput := &MealPlanGroceryListItemUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanGroceryListItemID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanGroceryListItemID := fake.BuildFakeID()

		exampleInput := &MealPlanGroceryListItemUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanGroceryListItem(ctx, "", mealPlanGroceryListItemID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanGroceryListItem ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanGroceryListItemUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanGroceryListItemID := fake.BuildFakeID()

		exampleInput := &MealPlanGroceryListItemUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanGroceryListItemID := fake.BuildFakeID()

		exampleInput := &MealPlanGroceryListItemUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanGroceryListItemID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID, exampleInput)

		assert.Error(t, err)
	})
}
