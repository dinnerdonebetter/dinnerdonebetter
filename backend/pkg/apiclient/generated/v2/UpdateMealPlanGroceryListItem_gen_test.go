// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_UpdateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/grocery_list_items/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanGroceryListItemID := fakes.BuildFakeID()

		data := fakes.BuildFakeMealPlanGroceryListItem()
		expected := &types.APIResponse[*types.MealPlanGroceryListItem]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanGroceryListItemID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanGroceryListItemID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanGroceryListItem(ctx, "", mealPlanGroceryListItemID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanGroceryListItem ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanGroceryListItemID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanGroceryListItemID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanGroceryListItemUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanGroceryListItemID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID, exampleInput)

		assert.Error(t, err)
	})
}
