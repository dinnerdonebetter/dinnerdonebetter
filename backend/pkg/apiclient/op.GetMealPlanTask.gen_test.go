// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMealPlanTask(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanTaskID := fake.BuildFakeID()

		data := &MealPlanTask{}
		expected := &APIResponse[*MealPlanTask]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanTaskID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMealPlanTask(ctx, mealPlanID, mealPlanTaskID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanTaskID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanTask(ctx, "", mealPlanTaskID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanTask ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanTask(ctx, mealPlanID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanTaskID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanTask(ctx, mealPlanID, mealPlanTaskID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanTaskID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanTaskID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanTask(ctx, mealPlanID, mealPlanTaskID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
