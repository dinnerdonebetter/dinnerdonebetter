// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_vessels/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()

		list := []*ValidVessel{}
		exampleResponse := &APIResponse[[]*ValidVessel]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[ValidVessel]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.SearchForValidVessels(ctx, q, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchForValidVessels(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		q := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchForValidVessels(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
