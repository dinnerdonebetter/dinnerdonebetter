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

func TestClient_CreateRecipePrepTask(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()

		data := fakes.BuildFakeRecipePrepTask()
		expected := &types.APIResponse[*types.RecipePrepTask]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.CreateRecipePrepTask(ctx, recipeID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.CreateRecipePrepTask(ctx, "", exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateRecipePrepTask(ctx, recipeID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.CreateRecipePrepTask(ctx, recipeID, exampleInput)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
