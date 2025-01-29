// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationInstrumentID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*ValidPreparationInstrument](t)

		expected := &APIResponse[*ValidPreparationInstrument]{
			Data: data,
		}

		exampleInput := fake.BuildFakeForTest[*ValidPreparationInstrumentUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateValidPreparationInstrument(ctx, validPreparationInstrumentID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid validPreparationInstrument ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fake.BuildFakeForTest[*ValidPreparationInstrumentUpdateRequestInput](t)

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateValidPreparationInstrument(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationInstrumentID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidPreparationInstrumentUpdateRequestInput](t)

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateValidPreparationInstrument(ctx, validPreparationInstrumentID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		validPreparationInstrumentID := fake.BuildFakeID()

		exampleInput := fake.BuildFakeForTest[*ValidPreparationInstrumentUpdateRequestInput](t)

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, validPreparationInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateValidPreparationInstrument(ctx, validPreparationInstrumentID, exampleInput)

		assert.Error(t, err)
	})
}
