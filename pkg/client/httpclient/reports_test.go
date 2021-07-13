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

func TestReports(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(reportsTestSuite))
}

type reportsBaseSuite struct {
	suite.Suite

	ctx           context.Context
	exampleReport *types.Report
}

var _ suite.SetupTestSuite = (*reportsBaseSuite)(nil)

func (s *reportsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleReport = fakes.BuildFakeReport()
}

type reportsTestSuite struct {
	suite.Suite

	reportsBaseSuite
}

func (s *reportsTestSuite) TestClient_ReportExists() {
	const expectedPathFormat = "/api/v1/reports/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleReport.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.ReportExists(s.ctx, s.exampleReport.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid report ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.ReportExists(s.ctx, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ReportExists(s.ctx, s.exampleReport.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.ReportExists(s.ctx, s.exampleReport.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *reportsTestSuite) TestClient_GetReport() {
	const expectedPathFormat = "/api/v1/reports/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleReport.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleReport)
		actual, err := c.GetReport(s.ctx, s.exampleReport.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleReport, actual)
	})

	s.Run("with invalid report ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetReport(s.ctx, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetReport(s.ctx, s.exampleReport.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleReport.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetReport(s.ctx, s.exampleReport.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *reportsTestSuite) TestClient_GetReports() {
	const expectedPath = "/api/v1/reports"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleReportList := fakes.BuildFakeReportList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleReportList)
		actual, err := c.GetReports(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleReportList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetReports(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetReports(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *reportsTestSuite) TestClient_CreateReport() {
	const expectedPath = "/api/v1/reports"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeReportCreationInput()
		exampleInput.BelongsToAccount = 0

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleReport)

		actual, err := c.CreateReport(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleReport, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateReport(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ReportCreationInput{}

		actual, err := c.CreateReport(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeReportCreationInputFromReport(s.exampleReport)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateReport(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeReportCreationInputFromReport(s.exampleReport)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateReport(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *reportsTestSuite) TestClient_UpdateReport() {
	const expectedPathFormat = "/api/v1/reports/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleReport.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleReport)

		err := c.UpdateReport(s.ctx, s.exampleReport)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateReport(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateReport(s.ctx, s.exampleReport)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateReport(s.ctx, s.exampleReport)
		assert.Error(t, err)
	})
}

func (s *reportsTestSuite) TestClient_ArchiveReport() {
	const expectedPathFormat = "/api/v1/reports/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleReport.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveReport(s.ctx, s.exampleReport.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid report ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveReport(s.ctx, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveReport(s.ctx, s.exampleReport.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveReport(s.ctx, s.exampleReport.ID)
		assert.Error(t, err)
	})
}

func (s *reportsTestSuite) TestClient_GetAuditLogForReport() {
	const (
		expectedPath   = "/api/v1/reports/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleReport.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForReport(s.ctx, s.exampleReport.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid report ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForReport(s.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForReport(s.ctx, s.exampleReport.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForReport(s.ctx, s.exampleReport.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
