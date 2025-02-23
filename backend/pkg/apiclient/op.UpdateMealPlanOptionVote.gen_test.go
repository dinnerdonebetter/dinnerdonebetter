// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s/votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		data := &MealPlanOptionVote{}
		expected := &APIResponse[*MealPlanOptionVote]{
			Data: data,
		}

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, "", mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, "", mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOption ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		mealPlanOptionVoteID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, "", mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOptionVote ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()
		mealPlanOptionID := fake.BuildFakeID()
		mealPlanOptionVoteID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionVoteUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)

		assert.Error(t, err)
	})
}
