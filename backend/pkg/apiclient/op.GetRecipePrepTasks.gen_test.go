// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRecipePrepTasks(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()

		list := []*RecipePrepTask{}
		exampleResponse := &APIResponse[[]*RecipePrepTask]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[RecipePrepTask]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, recipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.GetRecipePrepTasks(ctx, recipeID, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipePrepTasks(ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipePrepTasks(ctx, recipeID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		recipeID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, recipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipePrepTasks(ctx, recipeID, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
