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

func TestClient_SearchForValidIngredientStates(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_ingredient_states/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()

		list := fake.BuildFakeForTest[[]*ValidIngredientState](t)

		expected := &APIResponse[[]*ValidIngredientState]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.SearchForValidIngredientStates(ctx, q, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchForValidIngredientStates(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchForValidIngredientStates(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
