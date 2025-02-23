// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMealPlanOption(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()

		data := &MealPlanOption{}
		expected := &APIResponse[*MealPlanOption]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(ctx, "", mealPlanEventID, mealPlanOptionID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		mealPlanOptionID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(ctx, mealPlanID, "", mealPlanOptionID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOption ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(ctx, mealPlanID, mealPlanEventID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
