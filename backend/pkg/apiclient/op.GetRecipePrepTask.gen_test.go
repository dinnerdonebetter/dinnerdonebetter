// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/lib/internal/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRecipePrepTask(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipePrepTaskID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*RecipePrepTask](t)
		expected := &APIResponse[*RecipePrepTask]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipePrepTaskID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		recipePrepTaskID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipePrepTask(ctx, "", recipePrepTaskID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid recipePrepTask ID", func(t *testing.T) {
		t.Parallel()

		recipeID := fake.BuildFakeID()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipePrepTask(ctx, recipeID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipePrepTaskID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()
		recipePrepTaskID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, recipeID, recipePrepTaskID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipePrepTask(ctx, recipeID, recipePrepTaskID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
