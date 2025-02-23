// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s/options"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		data := &MealPlanOption{}
		expected := &APIResponse[*MealPlanOption]{
			Data: data,
		}

		exampleInput := &MealPlanOptionCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateMealPlanOption(ctx, mealPlanID, mealPlanEventID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid mealPlan ID", func(t *testing.T) {
		t.Parallel()

		mealPlanEventID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateMealPlanOption(ctx, "", mealPlanEventID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid mealPlanEvent ID", func(t *testing.T) {
		t.Parallel()

		mealPlanID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionCreationRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateMealPlanOption(ctx, mealPlanID, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionCreationRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateMealPlanOption(ctx, mealPlanID, mealPlanEventID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mealPlanID := fake.BuildFakeID()
		mealPlanEventID := fake.BuildFakeID()

		exampleInput := &MealPlanOptionCreationRequestInput{}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, mealPlanID, mealPlanEventID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateMealPlanOption(ctx, mealPlanID, mealPlanEventID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
