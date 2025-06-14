// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/accounts/instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInstrumentOwnershipID := fake.BuildFakeID()

		data := &AccountInstrumentOwnership{}
		expected := &APIResponse[*AccountInstrumentOwnership]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, accountInstrumentOwnershipID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetAccountInstrumentOwnership(ctx, accountInstrumentOwnershipID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid accountInstrumentOwnership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetAccountInstrumentOwnership(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInstrumentOwnershipID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetAccountInstrumentOwnership(ctx, accountInstrumentOwnershipID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInstrumentOwnershipID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, accountInstrumentOwnershipID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetAccountInstrumentOwnership(ctx, accountInstrumentOwnershipID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
