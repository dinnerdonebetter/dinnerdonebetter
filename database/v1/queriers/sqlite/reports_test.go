package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowFromReport(x *models.Report) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(reportsTableColumns).AddRow(
		x.ID,
		x.ReportType,
		x.Concern,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromReport(x *models.Report) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(reportsTableColumns).AddRow(
		x.ArchivedOn,
		x.ReportType,
		x.Concern,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestSqlite_buildGetReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleReportID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildGetReportQuery(exampleReportID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleReportID, args[1].(uint64))
	})
}

func TestSqlite_GetReport(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE belongs_to = ? AND id = ?"
		expected := &models.Report{
			ID: 123,
		}
		expectedUserID := uint64(321)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromReport(expected))

		actual, err := s.GetReport(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE belongs_to = ? AND id = ?"
		expected := &models.Report{
			ID: 123,
		}
		expectedUserID := uint64(321)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetReport(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetReportCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := s.buildGetReportCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetReportCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetReportCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllReportsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM reports WHERE archived_on IS NULL"

		actualQuery := s.buildGetAllReportsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllReportsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM reports WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllReportsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetReportsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := s.buildGetReportsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetReports(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM reports WHERE archived_on IS NULL"
		expectedReport := &models.Report{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.ReportList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			Reports: []models.Report{
				*expectedReport,
			},
		}

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromReport(expectedReport))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := s.GetReports(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetReports(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetReports(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning report", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.Report{
			ID: 321,
		}
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromReport(expected))

		actual, err := s.GetReports(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.Report{
			ID: 321,
		}
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM reports WHERE archived_on IS NULL"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromReport(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetReports(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_GetAllReportsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedReport := &models.Report{
			ID: 321,
		}
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromReport(expectedReport))

		expected := []models.Report{*expectedReport}
		actual, err := s.GetAllReportsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllReportsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllReportsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleReport := &models.Report{
			ID: 321,
		}
		expectedListQuery := "SELECT id, report_type, concern, created_on, updated_on, archived_on, belongs_to FROM reports WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromReport(exampleReport))

		actual, err := s.GetAllReportsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreateReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Report{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 3
		expectedQuery := "INSERT INTO reports (report_type,concern,belongs_to) VALUES (?,?,?)"
		actualQuery, args := s.buildCreateReportQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.ReportType, args[0].(string))
		assert.Equal(t, expected.Concern, args[1].(string))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
	})
}

func TestSqlite_CreateReport(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Report{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.ReportCreationInput{
			ReportType: expected.ReportType,
			Concern:    expected.Concern,
			BelongsTo:  expected.BelongsTo,
		}

		s, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO reports (report_type,concern,belongs_to) VALUES (?,?,?)"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.ReportType,
				expected.Concern,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM reports WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := s.CreateReport(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Report{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.ReportCreationInput{
			ReportType: expected.ReportType,
			Concern:    expected.Concern,
			BelongsTo:  expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO reports (report_type,concern,belongs_to) VALUES (?,?,?)"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.ReportType,
				expected.Concern,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := s.CreateReport(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdateReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Report{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 4
		expectedQuery := "UPDATE reports SET report_type = ?, concern = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildUpdateReportQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.ReportType, args[0].(string))
		assert.Equal(t, expected.Concern, args[1].(string))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
		assert.Equal(t, expected.ID, args[3].(uint64))
	})
}

func TestSqlite_UpdateReport(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Report{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE reports SET report_type = ?, concern = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.ReportType,
				expected.Concern,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := s.UpdateReport(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Report{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE reports SET report_type = ?, concern = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.ReportType,
				expected.Concern,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := s.UpdateReport(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchiveReportQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Report{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE reports SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := s.buildArchiveReportQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestSqlite_ArchiveReport(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Report{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE reports SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.ArchiveReport(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Report{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE reports SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := s.ArchiveReport(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
