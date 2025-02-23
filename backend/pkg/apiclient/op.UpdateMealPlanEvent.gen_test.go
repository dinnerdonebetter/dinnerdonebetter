// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		data := &MealPlanEvent{}
		expected := &APIResponse[*MealPlanEvent]{
			Data: data,
		}

		exampleInput := &MealPlanEventUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, mealPlanEventID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fake.BuildFakeID()

		exampleInput := &MealPlanEventUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanEvent(ctx, "", mealPlanEventID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanEventUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		exampleInput := &MealPlanEventUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, mealPlanEventID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		exampleInput := &MealPlanEventUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlanEvent(ctx, mealPlanID, mealPlanEventID, exampleInput)

		assert.Error(t, err)
	})
}
