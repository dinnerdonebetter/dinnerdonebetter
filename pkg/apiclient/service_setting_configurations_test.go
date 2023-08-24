package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestServiceSettingConfigurations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(serviceSettingConfigurationsTestSuite))
}

type serviceSettingConfigurationsBaseSuite struct {
	suite.Suite

	ctx                                context.Context
	exampleServiceSettingConfiguration *types.ServiceSettingConfiguration
}

var _ suite.SetupTestSuite = (*serviceSettingConfigurationsBaseSuite)(nil)

func (s *serviceSettingConfigurationsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleServiceSettingConfiguration = fakes.BuildFakeServiceSettingConfiguration()
}

type serviceSettingConfigurationsTestSuite struct {
	suite.Suite

	serviceSettingConfigurationsBaseSuite
}

func (s *serviceSettingConfigurationsTestSuite) TestClient_GetServiceSettingConfigurationForUserByName() {
	const expectedPath = "/api/v1/settings/configurations/user/%s"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)
		settingName := fakes.BuildFakeServiceSetting().Name

		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, settingName)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleServiceSettingConfiguration)
		actual, err := c.GetServiceSettingConfigurationForUserByName(s.ctx, settingName, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfiguration, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)
		settingName := fakes.BuildFakeServiceSetting().Name

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSettingConfigurationForUserByName(s.ctx, settingName, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)
		settingName := fakes.BuildFakeServiceSetting().Name

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, settingName)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSettingConfigurationForUserByName(s.ctx, settingName, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingConfigurationsTestSuite) TestClient_GetServiceSettingConfigurationsForUser() {
	const expectedPath = "/api/v1/settings/configurations/user"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleServiceSettingConfigurationList := fakes.BuildFakeServiceSettingConfigurationList()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleServiceSettingConfigurationList)
		actual, err := c.GetServiceSettingConfigurationsForUser(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfigurationList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSettingConfigurationsForUser(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSettingConfigurationsForUser(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingConfigurationsTestSuite) TestClient_GetServiceSettingConfigurationsForHousehold() {
	const expectedPath = "/api/v1/settings/configurations/household"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleServiceSettingConfigurationList := fakes.BuildFakeServiceSettingConfigurationList()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleServiceSettingConfigurationList)
		actual, err := c.GetServiceSettingConfigurationsForHousehold(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfigurationList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSettingConfigurationsForHousehold(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSettingConfigurationsForHousehold(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingConfigurationsTestSuite) TestClient_CreateServiceSettingConfiguration() {
	const expectedPath = "/api/v1/settings/configurations"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeServiceSettingConfigurationCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleServiceSettingConfiguration)

		actual, err := c.CreateServiceSettingConfiguration(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleServiceSettingConfiguration, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateServiceSettingConfiguration(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ServiceSettingConfigurationCreationRequestInput{}

		actual, err := c.CreateServiceSettingConfiguration(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(s.exampleServiceSettingConfiguration)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateServiceSettingConfiguration(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(s.exampleServiceSettingConfiguration)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateServiceSettingConfiguration(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingConfigurationsTestSuite) TestClient_UpdateServiceSettingConfiguration() {
	const expectedPathFormat = "/api/v1/settings/configurations/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleServiceSettingConfiguration.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleServiceSettingConfiguration)

		err := c.UpdateServiceSettingConfiguration(s.ctx, s.exampleServiceSettingConfiguration)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateServiceSettingConfiguration(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateServiceSettingConfiguration(s.ctx, s.exampleServiceSettingConfiguration)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateServiceSettingConfiguration(s.ctx, s.exampleServiceSettingConfiguration)
		assert.Error(t, err)
	})
}

func (s *serviceSettingConfigurationsTestSuite) TestClient_ArchiveServiceSettingConfiguration() {
	const expectedPathFormat = "/api/v1/settings/configurations/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleServiceSettingConfiguration.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveServiceSettingConfiguration(s.ctx, s.exampleServiceSettingConfiguration.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid service setting ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveServiceSettingConfiguration(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveServiceSettingConfiguration(s.ctx, s.exampleServiceSettingConfiguration.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveServiceSettingConfiguration(s.ctx, s.exampleServiceSettingConfiguration.ID)
		assert.Error(t, err)
	})
}
