package integration

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

func checkReportEquality(t *testing.T, expected, actual *models.Report) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ReportType, actual.ReportType, "expected ReportType for ID %d to be %v, but it was %v ", expected.ID, expected.ReportType, actual.ReportType)
	assert.Equal(t, expected.Concern, actual.Concern, "expected Concern for ID %d to be %v, but it was %v ", expected.ID, expected.Concern, actual.Concern)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyReport(t *testing.T) *models.Report {
	t.Helper()

	x := &models.ReportCreationInput{
		ReportType: fake.Word(),
		Concern:    fake.Word(),
	}
	y, err := todoClient.CreateReport(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestReports(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create report
			expected := &models.Report{
				ReportType: fake.Word(),
				Concern:    fake.Word(),
			}
			premade, err := todoClient.CreateReport(ctx, &models.ReportCreationInput{
				ReportType: expected.ReportType,
				Concern:    expected.Concern,
			})
			checkValueAndError(t, premade, err)

			// Assert report equality
			checkReportEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveReport(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetReport(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkReportEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create reports
			var expected []*models.Report
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyReport(t))
			}

			// Assert report list equality
			actual, err := todoClient.GetReports(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Reports),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Reports),
			)

			// Clean up
			for _, x := range actual.Reports {
				err = todoClient.ArchiveReport(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch report
			_, err := todoClient.GetReport(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create report
			expected := &models.Report{
				ReportType: fake.Word(),
				Concern:    fake.Word(),
			}
			premade, err := todoClient.CreateReport(ctx, &models.ReportCreationInput{
				ReportType: expected.ReportType,
				Concern:    expected.Concern,
			})
			checkValueAndError(t, premade, err)

			// Fetch report
			actual, err := todoClient.GetReport(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert report equality
			checkReportEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveReport(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateReport(ctx, &models.Report{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create report
			expected := &models.Report{
				ReportType: fake.Word(),
				Concern:    fake.Word(),
			}
			premade, err := todoClient.CreateReport(tctx, &models.ReportCreationInput{
				ReportType: fake.Word(),
				Concern:    fake.Word(),
			})
			checkValueAndError(t, premade, err)

			// Change report
			premade.Update(expected.ToInput())
			err = todoClient.UpdateReport(ctx, premade)
			assert.NoError(t, err)

			// Fetch report
			actual, err := todoClient.GetReport(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert report equality
			checkReportEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveReport(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create report
			expected := &models.Report{
				ReportType: fake.Word(),
				Concern:    fake.Word(),
			}
			premade, err := todoClient.CreateReport(ctx, &models.ReportCreationInput{
				ReportType: expected.ReportType,
				Concern:    expected.Concern,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveReport(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
