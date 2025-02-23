// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_UpdateServiceSettingConfiguration(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/settings/configurations/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		serviceSettingConfigurationID := fake.BuildFakeID()

		data := &ServiceSettingConfiguration{}
		expected := &APIResponse[*ServiceSettingConfiguration]{
			Data: data,
		}

		exampleInput := &ServiceSettingConfigurationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, serviceSettingConfigurationID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		err := c.UpdateServiceSettingConfiguration(ctx, serviceSettingConfigurationID, exampleInput)

		assert.NoError(t, err)

	})

	T.Run("with invalid serviceSettingConfiguration ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := &ServiceSettingConfigurationUpdateRequestInput{}

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		err := c.UpdateServiceSettingConfiguration(ctx, "", exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		serviceSettingConfigurationID := fake.BuildFakeID()

		exampleInput := &ServiceSettingConfigurationUpdateRequestInput{}

		c := buildTestClientWithInvalidURL(t)
		err := c.UpdateServiceSettingConfiguration(ctx, serviceSettingConfigurationID, exampleInput)

		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		serviceSettingConfigurationID := fake.BuildFakeID()

		exampleInput := &ServiceSettingConfigurationUpdateRequestInput{}

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, serviceSettingConfigurationID)
		c := buildTestClientWithInvalidResponse(t, spec)
		err := c.UpdateServiceSettingConfiguration(ctx, serviceSettingConfigurationID, exampleInput)

		assert.Error(t, err)
	})
}
