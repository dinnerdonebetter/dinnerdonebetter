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

func TestClient_UpdateMealPlan(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()

		data := fakes.BuildFakeMealPlan()
		expected := &types.APIResponse[*types.MealPlan]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeMealPlanUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateMealPlan(ctx, mealPlanID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeMealPlanUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateMealPlan(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateMealPlan(ctx, mealPlanID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mealPlanID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeMealPlanUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, mealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateMealPlan(ctx, mealPlanID, exampleInput)

		assert.Error(t, err)
	})
}
