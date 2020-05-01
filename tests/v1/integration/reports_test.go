package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkReportEquality(t *testing.T, expected, actual *models.Report) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ReportType, actual.ReportType, "expected ReportType for ID %d to be %v, but it was %v ", expected.ID, expected.ReportType, actual.ReportType)
	assert.Equal(t, expected.Concern, actual.Concern, "expected Concern for ID %d to be %v, but it was %v ", expected.ID, expected.Concern, actual.Concern)
	assert.NotZero(t, actual.CreatedOn)
}

func TestReports(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create report.
			exampleReport := fakemodels.BuildFakeReport()
			exampleReportInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := prixfixeClient.CreateReport(ctx, exampleReportInput)
			checkValueAndError(t, createdReport, err)

			// Assert report equality.
			checkReportEquality(t, exampleReport, createdReport)

			// Clean up.
			err = prixfixeClient.ArchiveReport(ctx, createdReport.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetReport(ctx, createdReport.ID)
			checkValueAndError(t, actual, err)
			checkReportEquality(t, exampleReport, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create reports.
			var expected []*models.Report
			for i := 0; i < 5; i++ {
				// Create report.
				exampleReport := fakemodels.BuildFakeReport()
				exampleReportInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)
				createdReport, reportCreationErr := prixfixeClient.CreateReport(ctx, exampleReportInput)
				checkValueAndError(t, createdReport, reportCreationErr)

				expected = append(expected, createdReport)
			}

			// Assert report list equality.
			actual, err := prixfixeClient.GetReports(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Reports),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Reports),
			)

			// Clean up.
			for _, createdReport := range actual.Reports {
				err = prixfixeClient.ArchiveReport(ctx, createdReport.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent report.
			actual, err := prixfixeClient.ReportExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant report exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create report.
			exampleReport := fakemodels.BuildFakeReport()
			exampleReportInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := prixfixeClient.CreateReport(ctx, exampleReportInput)
			checkValueAndError(t, createdReport, err)

			// Fetch report.
			actual, err := prixfixeClient.ReportExists(ctx, createdReport.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up report.
			assert.NoError(t, prixfixeClient.ArchiveReport(ctx, createdReport.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent report.
			_, err := prixfixeClient.GetReport(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create report.
			exampleReport := fakemodels.BuildFakeReport()
			exampleReportInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := prixfixeClient.CreateReport(ctx, exampleReportInput)
			checkValueAndError(t, createdReport, err)

			// Fetch report.
			actual, err := prixfixeClient.GetReport(ctx, createdReport.ID)
			checkValueAndError(t, actual, err)

			// Assert report equality.
			checkReportEquality(t, exampleReport, actual)

			// Clean up report.
			assert.NoError(t, prixfixeClient.ArchiveReport(ctx, createdReport.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleReport := fakemodels.BuildFakeReport()
			exampleReport.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateReport(ctx, exampleReport))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create report.
			exampleReport := fakemodels.BuildFakeReport()
			exampleReportInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := prixfixeClient.CreateReport(ctx, exampleReportInput)
			checkValueAndError(t, createdReport, err)

			// Change report.
			createdReport.Update(exampleReport.ToUpdateInput())
			err = prixfixeClient.UpdateReport(ctx, createdReport)
			assert.NoError(t, err)

			// Fetch report.
			actual, err := prixfixeClient.GetReport(ctx, createdReport.ID)
			checkValueAndError(t, actual, err)

			// Assert report equality.
			checkReportEquality(t, exampleReport, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up report.
			assert.NoError(t, prixfixeClient.ArchiveReport(ctx, createdReport.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveReport(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create report.
			exampleReport := fakemodels.BuildFakeReport()
			exampleReportInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)
			createdReport, err := prixfixeClient.CreateReport(ctx, exampleReportInput)
			checkValueAndError(t, createdReport, err)

			// Clean up report.
			assert.NoError(t, prixfixeClient.ArchiveReport(ctx, createdReport.ID))
		})
	})
}
