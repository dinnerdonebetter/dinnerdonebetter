package integration

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkReportEquality(t *testing.T, expected, actual *types.Report) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ReportType, actual.ReportType, "expected ReportType for report #%d to be %v, but it was %v ", expected.ID, expected.ReportType, actual.ReportType)
	assert.Equal(t, expected.Concern, actual.Concern, "expected Concern for report #%d to be %v, but it was %v ", expected.ID, expected.Concern, actual.Concern)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestReports_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create report.
			exampleReport := fakes.BuildFakeReport()
			exampleReportInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := testClients.main.CreateReport(ctx, exampleReportInput)
			requireNotNilAndNoProblems(t, createdReport, err)

			// assert report equality
			checkReportEquality(t, exampleReport, createdReport)

			auditLogEntries, err := testClients.admin.GetAuditLogForReport(ctx, createdReport.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ReportCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdReport.ID, audit.ReportAssignmentKey)

			// Clean up report.
			assert.NoError(t, testClients.main.ArchiveReport(ctx, createdReport.ID))
		}
	})
}

func (s *TestSuite) TestReports_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create reports
			var expected []*types.Report
			for i := 0; i < 5; i++ {
				exampleReport := fakes.BuildFakeReport()
				exampleReportInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)

				createdReport, reportCreationErr := testClients.main.CreateReport(ctx, exampleReportInput)
				requireNotNilAndNoProblems(t, createdReport, reportCreationErr)

				expected = append(expected, createdReport)
			}

			// assert report list equality
			actual, err := testClients.main.GetReports(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Reports),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Reports),
			)

			// clean up
			for _, createdReport := range actual.Reports {
				assert.NoError(t, testClients.main.ArchiveReport(ctx, createdReport.ID))
			}
		}
	})
}

func (s *TestSuite) TestReports_ExistenceChecking_ReturnsFalseForNonexistentReport() {
	s.runForEachClientExcept("should not return an error for nonexistent report", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.ReportExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestReports_ExistenceChecking_ReturnsTrueForValidReport() {
	s.runForEachClientExcept("should not return an error for existent report", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create report
			exampleReport := fakes.BuildFakeReport()
			exampleReportInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := testClients.main.CreateReport(ctx, exampleReportInput)
			requireNotNilAndNoProblems(t, createdReport, err)

			// retrieve report
			actual, err := testClients.main.ReportExists(ctx, createdReport.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up report
			assert.NoError(t, testClients.main.ArchiveReport(ctx, createdReport.ID))
		}
	})
}

func (s *TestSuite) TestReports_Reading_Returns404ForNonexistentReport() {
	s.runForEachClientExcept("it should return an error when trying to read a report that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetReport(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestReports_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create report
			exampleReport := fakes.BuildFakeReport()
			exampleReportInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := testClients.main.CreateReport(ctx, exampleReportInput)
			requireNotNilAndNoProblems(t, createdReport, err)

			// retrieve report
			actual, err := testClients.main.GetReport(ctx, createdReport.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert report equality
			checkReportEquality(t, exampleReport, actual)

			// clean up report
			assert.NoError(t, testClients.main.ArchiveReport(ctx, createdReport.ID))
		}
	})
}

func (s *TestSuite) TestReports_Updating_Returns404ForNonexistentReport() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleReport := fakes.BuildFakeReport()
			exampleReport.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateReport(ctx, exampleReport))
		}
	})
}

// convertReportToReportUpdateInput creates an ReportUpdateInput struct from a report.
func convertReportToReportUpdateInput(x *types.Report) *types.ReportUpdateInput {
	return &types.ReportUpdateInput{
		ReportType: x.ReportType,
		Concern:    x.Concern,
	}
}

func (s *TestSuite) TestReports_Updating() {
	s.runForEachClientExcept("it should be possible to update a report", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create report
			exampleReport := fakes.BuildFakeReport()
			exampleReportInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := testClients.main.CreateReport(ctx, exampleReportInput)
			requireNotNilAndNoProblems(t, createdReport, err)

			// change report
			createdReport.Update(convertReportToReportUpdateInput(exampleReport))
			assert.NoError(t, testClients.main.UpdateReport(ctx, createdReport))

			// retrieve changed report
			actual, err := testClients.main.GetReport(ctx, createdReport.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert report equality
			checkReportEquality(t, exampleReport, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForReport(ctx, createdReport.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ReportCreationEvent},
				{EventType: audit.ReportUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdReport.ID, audit.ReportAssignmentKey)

			// clean up report
			assert.NoError(t, testClients.main.ArchiveReport(ctx, createdReport.ID))
		}
	})
}

func (s *TestSuite) TestReports_Archiving_Returns404ForNonexistentReport() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveReport(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestReports_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a report", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create report
			exampleReport := fakes.BuildFakeReport()
			exampleReportInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := testClients.main.CreateReport(ctx, exampleReportInput)
			requireNotNilAndNoProblems(t, createdReport, err)

			// clean up report
			assert.NoError(t, testClients.main.ArchiveReport(ctx, createdReport.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForReport(ctx, createdReport.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ReportCreationEvent},
				{EventType: audit.ReportArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdReport.ID, audit.ReportAssignmentKey)
		}
	})
}

func (s *TestSuite) TestReports_Auditing_Returns404ForNonexistentReport() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForReport(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
