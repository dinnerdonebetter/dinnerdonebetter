package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromReports(reports ...*models.Report) *sqlmock.Rows {
	columns := reportsTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range reports {
		rowValues := []driver.Value{
			x.ID,
			x.ReportType,
			x.Concern,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromReport(x *models.Report) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(reportsTableColumns).AddRow(
		x.ArchivedOn,
		x.ReportType,
		x.Concern,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanReports(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanReports(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanReports(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildReportExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT reports.id FROM reports WHERE reports.id = $1 )"
		expectedArgs := []interface{}{
			exampleReport.ID,
		}
		actualQuery, actualArgs := p.buildReportExistsQuery(exampleReport.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ReportExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT reports.id FROM reports WHERE reports.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ReportExists(ctx, exampleReport.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ReportExists(ctx, exampleReport.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.id = $1"
		expectedArgs := []interface{}{
			exampleReport.ID,
		}
		actualQuery, actualArgs := p.buildGetReportQuery(exampleReport.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetReport(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ID,
			).
			WillReturnRows(buildMockRowsFromReports(exampleReport))

		actual, err := p.GetReport(ctx, exampleReport.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleReport, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetReport(ctx, exampleReport.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllReportsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL"
		actualQuery := p.buildGetAllReportsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllReportsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllReportsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfReportsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.id > $1 AND reports.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfReportsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllReports(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL"
	expectedGetQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.id > $1 AND reports.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleReportList := fakemodels.BuildFakeReportList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromReports(
					&exampleReportList.Reports[0],
					&exampleReportList.Reports[1],
					&exampleReportList.Reports[2],
				),
			)

		out := make(chan []models.Report)
		doneChan := make(chan bool, 1)

		err := p.GetAllReports(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.Report)

		err := p.GetAllReports(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []models.Report)

		err := p.GetAllReports(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.Report)

		err := p.GetAllReports(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleReport := fakemodels.BuildFakeReport()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromReport(exampleReport))

		out := make(chan []models.Report)

		err := p.GetAllReports(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetReportsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.archived_on IS NULL AND reports.created_on > $1 AND reports.created_on < $2 AND reports.last_updated_on > $3 AND reports.last_updated_on < $4 ORDER BY reports.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetReportsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetReports(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.archived_on IS NULL ORDER BY reports.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleReportList := fakemodels.BuildFakeReportList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromReports(
					&exampleReportList.Reports[0],
					&exampleReportList.Reports[1],
					&exampleReportList.Reports[2],
				),
			)

		actual, err := p.GetReports(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleReportList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetReports(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		actual, err := p.GetReports(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning report", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleReport := fakemodels.BuildFakeReport()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildErroneousMockRowFromReport(exampleReport),
			)

		actual, err := p.GetReports(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetReportsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM (SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS reports WHERE reports.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetReportsWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetReportsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleReportList := fakemodels.BuildFakeReportList()
		var exampleIDs []uint64
		for _, report := range exampleReportList.Reports {
			exampleIDs = append(exampleIDs, report.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM (SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS reports WHERE reports.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromReports(
					&exampleReportList.Reports[0],
					&exampleReportList.Reports[1],
					&exampleReportList.Reports[2],
				),
			)

		actual, err := p.GetReportsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleReportList.Reports, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM (SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS reports WHERE reports.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetReportsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM (SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS reports WHERE reports.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetReportsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning report", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM (SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.last_updated_on, reports.archived_on, reports.belongs_to_user FROM reports JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS reports WHERE reports.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleReport := fakemodels.BuildFakeReport()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromReport(exampleReport))

		actual, err := p.GetReportsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO reports (report_type,concern,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleReport.ReportType,
			exampleReport.Concern,
			exampleReport.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateReportQuery(exampleReport)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateReport(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO reports (report_type,concern,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleReport.ID, exampleReport.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleReport.ReportType,
				exampleReport.Concern,
				exampleReport.BelongsToUser,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateReport(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleReport, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleReport.ReportType,
				exampleReport.Concern,
				exampleReport.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateReport(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE reports SET report_type = $1, concern = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleReport.ReportType,
			exampleReport.Concern,
			exampleReport.BelongsToUser,
			exampleReport.ID,
		}
		actualQuery, actualArgs := p.buildUpdateReportQuery(exampleReport)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateReport(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE reports SET report_type = $1, concern = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ReportType,
				exampleReport.Concern,
				exampleReport.BelongsToUser,
				exampleReport.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateReport(ctx, exampleReport)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ReportType,
				exampleReport.Concern,
				exampleReport.BelongsToUser,
				exampleReport.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateReport(ctx, exampleReport)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE reports SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleReport.ID,
		}
		actualQuery, actualArgs := p.buildArchiveReportQuery(exampleReport.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveReport(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE reports SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleReport.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveReport(ctx, exampleReport.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleReport.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveReport(ctx, exampleReport.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleReport.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveReport(ctx, exampleReport.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
