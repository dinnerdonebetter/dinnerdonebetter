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

func TestClient_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()

		data := fakes.BuildFakeMealPlanEvent()
		expected := &types.APIResponse[*types.MealPlanEvent]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, mealPlanEventID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanEvent(ctx, "", mealPlanEventID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, mealPlanEventID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanEventID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, mealPlanEventID, exampleInput)

		assert.Error(t, err)
	})
}
