// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationInstrumentID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidPreparationInstrument()
		expected := &types.APIResponse[*types.ValidPreparationInstrument]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidPreparationInstrument(ctx, validPreparationInstrumentID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validPreparationInstrument ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidPreparationInstrument(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidPreparationInstrument(ctx, validPreparationInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidPreparationInstrument(ctx, validPreparationInstrumentID, exampleInput)

		assert.Error(t, err)
	})
}
