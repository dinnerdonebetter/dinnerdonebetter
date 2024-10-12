// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"fmt"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_SearchForUsers(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/users/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fakes.BuildFakeID()

		list := fakes.BuildFakeUsersList()
		for i := range list.Data {
			// the hashed passwords is never transmitted over the wire.
			list.Data[i].HashedPassword = ""
			// the two factor secret is transmitted over the wire only on creation.
			list.Data[i].TwoFactorSecret = ""
			// the two factor secret validation is never transmitted over the wire.
			list.Data[i].TwoFactorSecretVerifiedAt = nil
		}

		expected := &types.APIResponse[[]*types.User]{
			Pagination: &list.Pagination,
			Data:       list.Data,
		}

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.SearchForUsers(ctx, q, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, list, actual)
	})

	T.Run("with invalid query ", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchForUsers(ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fakes.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchForUsers(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=50&page=1&q=%s&sortBy=asc", q), expectedPathFormat)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchForUsers(ctx, q, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
