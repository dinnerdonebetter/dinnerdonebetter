package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestServiceSettings(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(serviceSettingsTestSuite))
}

type serviceSettingsBaseSuite struct {
	suite.Suite

	ctx                   context.Context
	exampleServiceSetting *types.ServiceSetting
}

var _ suite.SetupTestSuite = (*serviceSettingsBaseSuite)(nil)

func (s *serviceSettingsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleServiceSetting = fakes.BuildFakeServiceSetting()
}

type serviceSettingsTestSuite struct {
	suite.Suite

	serviceSettingsBaseSuite
}

func (s *serviceSettingsTestSuite) TestClient_GetServiceSetting() {
	const expectedPathFormat = "/api/v1/settings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleServiceSetting.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleServiceSetting)
		actual, err := c.GetServiceSetting(s.ctx, s.exampleServiceSetting.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleServiceSetting, actual)
	})

	s.Run("with invalid service setting ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetServiceSetting(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSetting(s.ctx, s.exampleServiceSetting.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleServiceSetting.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSetting(s.ctx, s.exampleServiceSetting.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingsTestSuite) TestClient_GetServiceSettings() {
	const expectedPath = "/api/v1/settings"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleServiceSettingList := fakes.BuildFakeServiceSettingList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleServiceSettingList)
		actual, err := c.GetServiceSettings(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetServiceSettings(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetServiceSettings(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingsTestSuite) TestClient_SearchServiceSettings() {
	const expectedPath = "/api/v1/settings/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		exampleServiceSettingList := fakes.BuildFakeServiceSettingList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleServiceSettingList.Data)
		actual, err := c.SearchServiceSettings(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingList.Data, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchServiceSettings(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchServiceSettings(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchServiceSettings(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingsTestSuite) TestClient_CreateServiceSetting() {
	const expectedPath = "/api/v1/settings"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeServiceSettingCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleServiceSetting)

		actual, err := c.CreateServiceSetting(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleServiceSetting, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateServiceSetting(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ServiceSettingCreationRequestInput{}

		actual, err := c.CreateServiceSetting(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertServiceSettingToServiceSettingCreationRequestInput(s.exampleServiceSetting)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateServiceSetting(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertServiceSettingToServiceSettingCreationRequestInput(s.exampleServiceSetting)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateServiceSetting(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *serviceSettingsTestSuite) TestClient_UpdateServiceSetting() {
	const expectedPathFormat = "/api/v1/settings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleServiceSetting.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleServiceSetting)

		err := c.UpdateServiceSetting(s.ctx, s.exampleServiceSetting)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateServiceSetting(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateServiceSetting(s.ctx, s.exampleServiceSetting)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateServiceSetting(s.ctx, s.exampleServiceSetting)
		assert.Error(t, err)
	})
}

func (s *serviceSettingsTestSuite) TestClient_ArchiveServiceSetting() {
	const expectedPathFormat = "/api/v1/settings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleServiceSetting.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveServiceSetting(s.ctx, s.exampleServiceSetting.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid service setting ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveServiceSetting(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveServiceSetting(s.ctx, s.exampleServiceSetting.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveServiceSetting(s.ctx, s.exampleServiceSetting.ID)
		assert.Error(t, err)
	})
}
