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

func TestClient_UpdateMealPlanOption(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		data := fakes.BuildFakeMealPlanOption()
		expected := &types.APIResponse[*types.MealPlanOption]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOption(ctx, "", mealPlanEventID, mealPlanOptionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()

		mealPlanOptionID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOption(ctx, mealPlanID, "", mealPlanOptionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanOption ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanOption(ctx, mealPlanID, mealPlanEventID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()
		mealPlanOptionID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID, mealPlanOptionID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)

		assert.Error(t, err)
	})
}
