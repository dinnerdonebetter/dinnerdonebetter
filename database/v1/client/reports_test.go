package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ReportExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("ReportExists", mock.Anything, exampleReport.ID).Return(true, nil)

		actual, err := c.ReportExists(ctx, exampleReport.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return(exampleReport, nil)

		actual, err := c.GetReport(ctx, exampleReport.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleReport, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllReportsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetAllReportsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllReportsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetReports(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleReportList := fakemodels.BuildFakeReportList()

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetReports", mock.Anything, filter).Return(exampleReportList, nil)

		actual, err := c.GetReports(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleReportList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleReportList := fakemodels.BuildFakeReportList()

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetReports", mock.Anything, filter).Return(exampleReportList, nil)

		actual, err := c.GetReports(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleReportList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("CreateReport", mock.Anything, exampleInput).Return(exampleReport, nil)

		actual, err := c.CreateReport(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleReport, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()

		mockDB.ReportDataManager.On("UpdateReport", mock.Anything, exampleReport).Return(expected)

		err := c.UpdateReport(ctx, exampleReport)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("ArchiveReport", mock.Anything, exampleReport.ID, exampleReport.BelongsToUser).Return(expected)

		err := c.ArchiveReport(ctx, exampleReport.ID, exampleReport.BelongsToUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
