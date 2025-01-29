// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/lib/internal/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*MealPlanOptionVote](t)
		expected := &APIResponse[*MealPlanOptionVote]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(ctx, "", mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(ctx, mealPlanID, "", mealPlanOptionID, mealPlanOptionVoteID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOption ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		mealPlanOptionVoteID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, "", mealPlanOptionVoteID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOptionVote ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
