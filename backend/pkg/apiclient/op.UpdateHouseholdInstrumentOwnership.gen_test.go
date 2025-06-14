// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateAccountInstrumentOwnership(T *testing.T) {
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

		exampleInput := &AccountInstrumentOwnershipUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, accountInstrumentOwnershipID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateAccountInstrumentOwnership(ctx, accountInstrumentOwnershipID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid accountInstrumentOwnership ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &AccountInstrumentOwnershipUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateAccountInstrumentOwnership(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInstrumentOwnershipID := fake.BuildFakeID()

		exampleInput := &AccountInstrumentOwnershipUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateAccountInstrumentOwnership(ctx, accountInstrumentOwnershipID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountInstrumentOwnershipID := fake.BuildFakeID()

		exampleInput := &AccountInstrumentOwnershipUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, accountInstrumentOwnershipID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateAccountInstrumentOwnership(ctx, accountInstrumentOwnershipID, exampleInput)

		assert.Error(t, err)
	})
}
