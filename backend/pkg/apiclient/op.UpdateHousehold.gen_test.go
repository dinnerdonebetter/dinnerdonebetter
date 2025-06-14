// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateAccount(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/accounts/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()

		data := &Account{}
		expected := &APIResponse[*Account]{
			Data: data,
		}

		exampleInput := &AccountUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, accountID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateAccount(ctx, accountID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid account ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &AccountUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateAccount(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()

		exampleInput := &AccountUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateAccount(ctx, accountID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		accountID := fake.BuildFakeID()

		exampleInput := &AccountUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, accountID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateAccount(ctx, accountID, exampleInput)

		assert.Error(t, err)
	})
}
