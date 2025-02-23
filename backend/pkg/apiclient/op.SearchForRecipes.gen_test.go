// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchForRecipes(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/recipes/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()

		list := []*Recipe{}
		exampleResponse := &APIResponse[[]*Recipe]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[Recipe]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.SearchForRecipes(ctx, q, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchForRecipes(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchForRecipes(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
