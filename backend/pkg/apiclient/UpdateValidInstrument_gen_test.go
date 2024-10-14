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

func TestClient_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validInstrumentID := fakes.BuildFakeID()

		data := fakes.BuildFakeValidInstrument()
		expected := &types.APIResponse[*types.ValidInstrument]{
			Data: data,
		}

		exampleInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidInstrument(ctx, validInstrumentID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validInstrument ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidInstrument(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidInstrument(ctx, validInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validInstrumentID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidInstrument(ctx, validInstrumentID, exampleInput)

		assert.Error(t, err)
	})
}
