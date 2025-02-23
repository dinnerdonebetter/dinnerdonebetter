// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetValidInstrument(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validInstrumentID := fake.BuildFakeID()

		data := &ValidInstrument{}
		expected := &APIResponse[*ValidInstrument]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validInstrumentID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetValidInstrument(ctx, validInstrumentID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid validInstrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidInstrument(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validInstrumentID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidInstrument(ctx, validInstrumentID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		validInstrumentID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, validInstrumentID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidInstrument(ctx, validInstrumentID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
