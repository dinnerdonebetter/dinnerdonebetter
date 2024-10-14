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

func TestClient_GetMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		list := fakes.BuildFakeMealPlanOptionVotesList()

		expected := &types.APIResponse[[]*types.MealPlanOptionVote]{
			Pagination: &list.Pagination,
			Data:       list.Data,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with empty mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVotes(ctx, "", mealPlanEventID, mealPlanOptionID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with empty mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()

		mealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVotes(ctx, mealPlanID, "", mealPlanOptionID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with empty mealPlanOption ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
