// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()

		data := &MealPlanEvent{}
		expected := &APIResponse[*MealPlanEvent]{
			Data: data,
		}

		exampleInput := &MealPlanEventCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, mealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateMealPlanEvent(ctx, mealPlanID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &MealPlanEventCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateMealPlanEvent(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanEventCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateMealPlanEvent(ctx, mealPlanID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanEventCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, mealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateMealPlanEvent(ctx, mealPlanID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
