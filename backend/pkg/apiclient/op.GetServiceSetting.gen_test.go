// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetServiceSetting(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/settings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		serviceSettingID := fake.BuildFakeID()

		data := &ServiceSetting{}
		expected := &APIResponse[*ServiceSetting]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, serviceSettingID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.GetServiceSetting(ctx, serviceSettingID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid serviceSetting ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetServiceSetting(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		serviceSettingID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSetting(ctx, serviceSettingID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		serviceSettingID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, serviceSettingID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSetting(ctx, serviceSettingID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
