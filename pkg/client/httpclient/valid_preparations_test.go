package httpclient

import (
	"context"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidPreparations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validPreparationsTestSuite))
}

type validPreparationsBaseSuite struct {
	suite.Suite

	ctx                     context.Context
	exampleValidPreparation *types.ValidPreparation
}

var _ suite.SetupTestSuite = (*validPreparationsBaseSuite)(nil)

func (s *validPreparationsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidPreparation = fakes.BuildFakeValidPreparation()
}

type validPreparationsTestSuite struct {
	suite.Suite

	validPreparationsBaseSuite
}

func (s *validPreparationsTestSuite) TestClient_ValidPreparationExists() {
	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleValidPreparation.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.ValidPreparationExists(s.ctx, s.exampleValidPreparation.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.ValidPreparationExists(s.ctx, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidPreparationExists(s.ctx, s.exampleValidPreparation.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.ValidPreparationExists(s.ctx, s.exampleValidPreparation.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *validPreparationsTestSuite) TestClient_GetValidPreparation() {
	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparation)
		actual, err := c.GetValidPreparation(s.ctx, s.exampleValidPreparation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparation, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparation(s.ctx, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparation(s.ctx, s.exampleValidPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparation(s.ctx, s.exampleValidPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationsTestSuite) TestClient_GetValidPreparations() {
	const expectedPath = "/api/v1/valid_preparations"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidPreparationList)
		actual, err := c.GetValidPreparations(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationsTestSuite) TestClient_SearchValidPreparations() {
	const expectedPath = "/api/v1/valid_preparations/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidPreparationList.ValidPreparations)
		actual, err := c.SearchValidPreparations(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList.ValidPreparations, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidPreparations(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidPreparations(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidPreparations(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationsTestSuite) TestClient_CreateValidPreparation() {
	const expectedPath = "/api/v1/valid_preparations"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparation)

		actual, err := c.CreateValidPreparation(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidPreparation, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidPreparation(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidPreparationCreationInput{}

		actual, err := c.CreateValidPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(s.exampleValidPreparation)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(s.exampleValidPreparation)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationsTestSuite) TestClient_UpdateValidPreparation() {
	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparation)

		err := c.UpdateValidPreparation(s.ctx, s.exampleValidPreparation)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidPreparation(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidPreparation(s.ctx, s.exampleValidPreparation)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidPreparation(s.ctx, s.exampleValidPreparation)
		assert.Error(t, err)
	})
}

func (s *validPreparationsTestSuite) TestClient_ArchiveValidPreparation() {
	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidPreparation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidPreparation(s.ctx, s.exampleValidPreparation.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidPreparation(s.ctx, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidPreparation(s.ctx, s.exampleValidPreparation.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidPreparation(s.ctx, s.exampleValidPreparation.ID)
		assert.Error(t, err)
	})
}

func (s *validPreparationsTestSuite) TestClient_GetAuditLogForValidPreparation() {
	const (
		expectedPath   = "/api/v1/valid_preparations/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleValidPreparation.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForValidPreparation(s.ctx, s.exampleValidPreparation.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForValidPreparation(s.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForValidPreparation(s.ctx, s.exampleValidPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForValidPreparation(s.ctx, s.exampleValidPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
