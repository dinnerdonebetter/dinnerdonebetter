// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetServiceSettingConfigurationByName(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/settings/configurations/user/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		serviceSettingConfigurationName := fake.BuildFakeID()

		list := []*ServiceSettingConfiguration{}
		exampleResponse := &APIResponse[[]*ServiceSettingConfiguration]{
			Pagination: fake.BuildFakeForTest[*Pagination](t),
			Data:       list,
		}
		expected := &QueryFilteredResult[ServiceSettingConfiguration]{
			Pagination: *exampleResponse.Pagination,
			Data:       list,
		}

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, serviceSettingConfigurationName)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)
		actual, err := c.GetServiceSettingConfigurationByName(ctx, serviceSettingConfigurationName, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	T.Run("with empty serviceSettingConfigurationName ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetServiceSettingConfigurationByName(ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		serviceSettingConfigurationName := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSettingConfigurationByName(ctx, serviceSettingConfigurationName, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		serviceSettingConfigurationName := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPathFormat, serviceSettingConfigurationName)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSettingConfigurationByName(ctx, serviceSettingConfigurationName, nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
