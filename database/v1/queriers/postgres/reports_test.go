package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromReport(reports ...*models.Report) *sqlmock.Rows {
	includeCount := len(reports) > 1
	columns := reportsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range reports {
		rowValues := []driver.Value{
			x.ID,
			x.ReportType,
			x.Concern,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		if includeCount {
			rowValues = append(rowValues, len(reports))
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
		x.UpdatedOn,
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

		_, _, err := p.scanReports(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanReports(mockRows)
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

		expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.id = $1"
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
	expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.updated_on, reports.archived_on, reports.belongs_to_user FROM reports WHERE reports.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleReport.ID,
			).
			WillReturnRows(buildMockRowsFromReport(exampleReport))

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

func TestPostgres_buildGetReportsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.updated_on, reports.archived_on, reports.belongs_to_user, (SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL) FROM reports WHERE reports.archived_on IS NULL AND reports.created_on > $1 AND reports.created_on < $2 AND reports.updated_on > $3 AND reports.updated_on < $4 ORDER BY reports.id LIMIT 20 OFFSET 180"
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

	expectedListQuery := "SELECT reports.id, reports.report_type, reports.concern, reports.created_on, reports.updated_on, reports.archived_on, reports.belongs_to_user, (SELECT COUNT(reports.id) FROM reports WHERE reports.archived_on IS NULL) FROM reports WHERE reports.archived_on IS NULL ORDER BY reports.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleReportList := fakemodels.BuildFakeReportList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(
				buildMockRowsFromReport(
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(
				buildErroneousMockRowFromReport(exampleReport),
			)

		actual, err := p.GetReports(ctx, filter)
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

		expectedQuery := "UPDATE reports SET report_type = $1, concern = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING updated_on"
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

	expectedQuery := "UPDATE reports SET report_type = $1, concern = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
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

		expectedQuery := "UPDATE reports SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"
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

	expectedQuery := "UPDATE reports SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"

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
