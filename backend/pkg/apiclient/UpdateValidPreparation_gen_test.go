// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidPreparation()
		expected := &types.APIResponse[*types.ValidPreparation]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidPreparationUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidPreparation(ctx, validPreparationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validPreparation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidPreparationUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidPreparation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidPreparationUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidPreparation(ctx, validPreparationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidPreparationUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidPreparation(ctx, validPreparationID, exampleInput)

		assert.Error(t, err)
	})
}
