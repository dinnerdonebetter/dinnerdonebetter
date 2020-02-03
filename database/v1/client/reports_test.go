package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleReportID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.Report{}

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetReport", mock.Anything, exampleReportID, exampleUserID).Return(expected, nil)

		actual, err := c.GetReport(context.Background(), exampleReportID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetReportCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetReportCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetReportCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetReportCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetReportCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllReportsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("GetAllReportsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllReportsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetReports(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.ReportList{}

		mockDB.ReportDataManager.On("GetReports", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetReports(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.ReportList{}

		mockDB.ReportDataManager.On("GetReports", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetReports(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.ReportCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.Report{}

		mockDB.ReportDataManager.On("CreateReport", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateReport(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.Report{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.ReportDataManager.On("UpdateReport", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateReport(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveReport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleReportID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.ReportDataManager.On("ArchiveReport", mock.Anything, exampleReportID, exampleUserID).Return(expected)

		err := c.ArchiveReport(context.Background(), exampleUserID, exampleReportID)
		assert.NoError(t, err)
	})
}
