// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateMealPlanTaskStatus(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanTaskID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*MealPlanTask](t)

		expected := &APIResponse[*MealPlanTask]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*MealPlanTaskStatusChangeRequestInput](t)

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, mealPlanID, mealPlanTaskID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, mealPlanTaskID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanTaskID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanTaskStatusChangeRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanTaskStatus(ctx, "", mealPlanTaskID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanTask ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanTaskStatusChangeRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanTaskID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanTaskStatusChangeRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, mealPlanTaskID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanTaskID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*MealPlanTaskStatusChangeRequestInput](t)

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, mealPlanID, mealPlanTaskID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanTaskStatus(ctx, mealPlanID, mealPlanTaskID, exampleInput)

		assert.Error(t, err)
	})
}
