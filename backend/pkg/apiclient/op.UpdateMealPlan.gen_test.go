// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateMealPlan(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()

		data := &MealPlan{}
		expected := &APIResponse[*MealPlan]{
			Data: data,
		}

		exampleInput := &MealPlanUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlan(ctx, mealPlanID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &MealPlanUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlan(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlan(ctx, mealPlanID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlan(ctx, mealPlanID, exampleInput)

		assert.Error(t, err)
	})
}
