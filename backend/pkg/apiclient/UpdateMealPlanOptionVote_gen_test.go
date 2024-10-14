// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()
		mealPlanOptionVoteID := fakes.BuildFakeID()

		data := fakes.BuildFakeMealPlanOptionVote()
		expected := &types.APIResponse[*types.MealPlanOptionVote]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()
		mealPlanOptionVoteID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, "", mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()

		mealPlanOptionID := fakes.BuildFakeID()
		mealPlanOptionVoteID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, "", mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOption ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()

		mealPlanOptionVoteID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, "", mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOptionVote ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()
		mealPlanOptionVoteID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()
		mealPlanOptionVoteID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})
}
