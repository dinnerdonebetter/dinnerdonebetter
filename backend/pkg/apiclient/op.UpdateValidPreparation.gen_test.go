// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationID := fake.BuildFakeID()

		data := &ValidPreparation{}
		expected := &APIResponse[*ValidPreparation]{
			Data: data,
		}

		exampleInput := &ValidPreparationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidPreparation(ctx, validPreparationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validPreparation ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &ValidPreparationUpdateRequestInput{}

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidPreparation(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationID := fake.BuildFakeID()

		exampleInput := &ValidPreparationUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidPreparation(ctx, validPreparationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validPreparationID := fake.BuildFakeID()

		exampleInput := &ValidPreparationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidPreparation(ctx, validPreparationID, exampleInput)

		assert.Error(t, err)
	})
}
