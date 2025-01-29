// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/vote"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		data := []*MealPlanOptionVote{fake.BuildFakeForTest[*MealPlanOptionVote](t)}

		expected := &APIResponse[[]*MealPlanOptionVote]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*MealPlanOptionVoteCreationRequestInput](t)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanOptionVoteCreationRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateMealPlanOptionVote(ctx, "", mealPlanEventID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanOptionVoteCreationRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateMealPlanOptionVote(ctx, mealPlanID, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanOptionVoteCreationRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanOptionVoteCreationRequestInput](t)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
