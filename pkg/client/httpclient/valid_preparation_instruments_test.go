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

func TestValidPreparationInstruments(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validPreparationInstrumentsTestSuite))
}

type validPreparationInstrumentsBaseSuite struct {
	suite.Suite

	ctx                               context.Context
	exampleValidPreparationInstrument *types.ValidPreparationInstrument
}

var _ suite.SetupTestSuite = (*validPreparationInstrumentsBaseSuite)(nil)

func (s *validPreparationInstrumentsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
}

type validPreparationInstrumentsTestSuite struct {
	suite.Suite

	validPreparationInstrumentsBaseSuite
}

func (s *validPreparationInstrumentsTestSuite) TestClient_ValidPreparationInstrumentExists() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.ValidPreparationInstrumentExists(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid valid preparation instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.ValidPreparationInstrumentExists(s.ctx, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidPreparationInstrumentExists(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.ValidPreparationInstrumentExists(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetValidPreparationInstrument() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrument)
		actual, err := c.GetValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationInstrument, actual)
	})

	s.Run("with invalid valid preparation instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationInstrument(s.ctx, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetValidPreparationInstruments() {
	const expectedPath = "/api/v1/valid_preparation_instruments"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidPreparationInstrumentList)
		actual, err := c.GetValidPreparationInstruments(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationInstruments(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationInstruments(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_CreateValidPreparationInstrument() {
	const expectedPath = "/api/v1/valid_preparation_instruments"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrument)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidPreparationInstrument, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidPreparationInstrumentCreationInput{}

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(s.exampleValidPreparationInstrument)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(s.exampleValidPreparationInstrument)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidPreparationInstrument(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_UpdateValidPreparationInstrument() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationInstrument)

		err := c.UpdateValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidPreparationInstrument(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_ArchiveValidPreparationInstrument() {
	const expectedPathFormat = "/api/v1/valid_preparation_instruments/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidPreparationInstrument.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid preparation instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidPreparationInstrument(s.ctx, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
	})
}

func (s *validPreparationInstrumentsTestSuite) TestClient_GetAuditLogForValidPreparationInstrument() {
	const (
		expectedPath   = "/api/v1/valid_preparation_instruments/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleValidPreparationInstrument.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid valid preparation instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForValidPreparationInstrument(s.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForValidPreparationInstrument(s.ctx, s.exampleValidPreparationInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
