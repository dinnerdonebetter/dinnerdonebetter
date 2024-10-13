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

func TestClient_UpdateMealPlanTaskStatus(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanTaskID := fakes.BuildFakeID()

		data := fakes.BuildFakeMealPlanTask()
		expected := &types.APIResponse[*types.MealPlanTask]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, mealPlanID, mealPlanTaskID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, mealPlanTaskID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanTaskID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanTaskStatus(ctx, "", mealPlanTaskID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanTask ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanTaskID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, mealPlanTaskID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()
		mealPlanTaskID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, mealPlanID, mealPlanTaskID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, mealPlanTaskID, exampleInput)

		assert.Error(t, err)
	})
}
