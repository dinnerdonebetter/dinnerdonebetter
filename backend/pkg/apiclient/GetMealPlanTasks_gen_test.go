// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMealPlanTasks(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()

		list := fakes.BuildFakeMealPlanTasksList()

		expected := &types.APIResponse[[]*types.MealPlanTask]{
			Pagination: &list.Pagination,
			Data:       list.Data,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, mealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMealPlanTasks(ctx, mealPlanID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with empty mealPlan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanTasks(ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanTasks(ctx, mealPlanID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, mealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanTasks(ctx, mealPlanID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
