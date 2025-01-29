// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validInstrumentID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidInstrument](t)

		expected := &APIResponse[*ValidInstrument]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidInstrumentUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidInstrument(ctx, validInstrumentID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validInstrument ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidInstrumentUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidInstrument(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validInstrumentID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidInstrumentUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidInstrument(ctx, validInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validInstrumentID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidInstrumentUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidInstrument(ctx, validInstrumentID, exampleInput)

		assert.Error(t, err)
	})
}
