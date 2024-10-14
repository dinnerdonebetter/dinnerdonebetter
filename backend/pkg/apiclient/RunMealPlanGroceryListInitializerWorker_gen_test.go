// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_RunMealPlanGroceryListInitializerWorker(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/workers/meal_plan_grocery_list_init"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		data := fakes.BuildFakeInitializeMealPlanGroceryListResponse()
		expected := &types.APIResponse[*types.InitializeMealPlanGroceryListResponse]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeInitializeMealPlanGroceryListRequest()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.RunMealPlanGroceryListInitializerWorker(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fakes.BuildFakeInitializeMealPlanGroceryListRequest()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RunMealPlanGroceryListInitializerWorker(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleInput := fakes.BuildFakeInitializeMealPlanGroceryListRequest()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.RunMealPlanGroceryListInitializerWorker(ctx, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
